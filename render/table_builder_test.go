package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTableBuilder(t *testing.T) {
	tb := NewTableBuilder("NAME", "TYPE", "STATUS")
	td := tb.Build()

	assert.Equal(t, []string{"NAME", "TYPE", "STATUS"}, td.Headers)
	assert.Empty(t, td.Rows)
}

func TestTableBuilder_AddRow(t *testing.T) {
	tb := NewTableBuilder("NAME", "VALUE")
	tb.AddRow("foo", "bar")
	tb.AddRow("baz", "qux")

	td := tb.Build()
	assert.Len(t, td.Rows, 2)
	assert.Equal(t, []string{"foo", "bar"}, td.Rows[0])
	assert.Equal(t, []string{"baz", "qux"}, td.Rows[1])
}

func TestTableBuilder_Fluent(t *testing.T) {
	td := NewTableBuilder("A", "B").
		AddRow("1", "2").
		AddRow("3", "4").
		Build()

	assert.Len(t, td.Headers, 2)
	assert.Len(t, td.Rows, 2)
}

func TestTableBuilder_Len(t *testing.T) {
	tb := NewTableBuilder("X")
	assert.Equal(t, 0, tb.Len())
	tb.AddRow("a")
	assert.Equal(t, 1, tb.Len())
	tb.AddRow("b")
	assert.Equal(t, 2, tb.Len())
}

func TestTableBuilder_EmptyHeaders(t *testing.T) {
	tb := NewTableBuilder()
	td := tb.Build()
	assert.Empty(t, td.Headers)
	assert.Empty(t, td.Rows)
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string unchanged", "hello", 10, "hello"},
		{"exact length unchanged", "hello", 5, "hello"},
		{"long string truncated", "hello world", 8, "hello..."},
		{"empty string", "", 10, ""},
		{"maxLen less than 4", "hello", 3, "hello"},
		{"maxLen exactly 4", "hello", 4, "h..."},
		{"one char over", "abcdef", 5, "ab..."},
		{"unicode safe", "description text here", 15, "description ..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.want, got)
		})
	}
}
