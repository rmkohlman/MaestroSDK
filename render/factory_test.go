package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRendererFactory(t *testing.T) {
	f := NewRendererFactory()
	assert.NotNil(t, f)

	// Should have all 6 built-in renderers
	names := f.Names()
	assert.Len(t, names, 6)
	assert.Contains(t, names, RendererJSON)
	assert.Contains(t, names, RendererYAML)
	assert.Contains(t, names, RendererPlain)
	assert.Contains(t, names, RendererColored)
	assert.Contains(t, names, RendererTable)
	assert.Contains(t, names, RendererCompact)
}

func TestRendererFactory_Create(t *testing.T) {
	f := NewRendererFactory()

	tests := []struct {
		name     RendererName
		wantName RendererName
		wantNil  bool
	}{
		{RendererJSON, RendererJSON, false},
		{RendererYAML, RendererYAML, false},
		{RendererPlain, RendererPlain, false},
		{RendererColored, RendererColored, false},
		{RendererTable, RendererTable, false},
		{RendererCompact, RendererCompact, false},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			r := f.Create(tt.name)
			if tt.wantNil {
				assert.Nil(t, r)
			} else {
				require.NotNil(t, r)
				assert.Equal(t, tt.wantName, r.Name())
			}
		})
	}
}

func TestRendererFactory_CreateReturnsNewInstances(t *testing.T) {
	f := NewRendererFactory()

	// Each call should return a new instance
	r1 := f.Create(RendererJSON)
	r2 := f.Create(RendererJSON)

	require.NotNil(t, r1)
	require.NotNil(t, r2)
	// They should be different pointers
	assert.NotSame(t, r1, r2)
}

func TestRendererFactory_RegisterAll(t *testing.T) {
	f := NewRendererFactory()
	f.RegisterAll()

	// Verify all renderers are in the global registry
	for _, name := range f.Names() {
		r := Get(name)
		assert.NotNil(t, r, "renderer %s should be registered", name)
		assert.Equal(t, name, r.Name())
	}
}

func TestRendererFactory_RegisterOne(t *testing.T) {
	f := NewRendererFactory()

	t.Run("known renderer", func(t *testing.T) {
		ok := f.RegisterOne(RendererJSON)
		assert.True(t, ok)

		r := Get(RendererJSON)
		assert.NotNil(t, r)
		assert.Equal(t, RendererJSON, r.Name())
	})

	t.Run("unknown renderer", func(t *testing.T) {
		ok := f.RegisterOne("nonexistent")
		assert.False(t, ok)
	})
}

func TestRendererFactory_AddCreator(t *testing.T) {
	f := NewRendererFactory()

	customName := RendererName("custom")
	f.AddCreator(customName, func() Renderer {
		return NewPlainRenderer() // just reuse plain for testing
	})

	r := f.Create(customName)
	require.NotNil(t, r)
	// The custom creator returns a PlainRenderer
	assert.Equal(t, RendererPlain, r.Name())

	// Should now appear in names
	assert.Contains(t, f.Names(), customName)
}

func TestRendererFactory_CreateColored(t *testing.T) {
	f := NewRendererFactory()

	t.Run("with default icons", func(t *testing.T) {
		r := f.CreateColored(DefaultIcons())
		assert.NotNil(t, r)
		assert.Equal(t, RendererColored, r.Name())
	})

	t.Run("with nerd font icons", func(t *testing.T) {
		r := f.CreateColored(NerdFontIcons())
		assert.NotNil(t, r)
		assert.Equal(t, RendererColored, r.Name())
	})

	t.Run("renders with custom icons", func(t *testing.T) {
		r := f.CreateColored(PlainIcons())
		var buf bytes.Buffer
		err := r.RenderMessage(&buf, Message{
			Level:   LevelSuccess,
			Content: "test",
		})
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "test")
	})
}

func TestFactoryCreatedRenderers_Functional(t *testing.T) {
	f := NewRendererFactory()

	// Verify each factory-created renderer actually works
	data := NewOrderedKeyValueData(
		KeyValue{Key: "key", Value: "value"},
	)
	opts := Options{Type: TypeKeyValue}

	for _, name := range f.Names() {
		t.Run(string(name), func(t *testing.T) {
			r := f.Create(name)
			require.NotNil(t, r)

			var buf bytes.Buffer
			err := r.Render(&buf, data, opts)
			require.NoError(t, err)

			// All renderers should produce some output for key-value data
			// (except table renderer which suppresses non-table types)
			if name == RendererTable {
				// Table renderer outputs key-value as simple pairs
				assert.Contains(t, buf.String(), "key")
			} else {
				assert.NotEmpty(t, buf.String())
			}
		})
	}
}
