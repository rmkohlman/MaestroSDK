package render

// RendererFactory centralizes the creation and registration of all
// renderer implementations. It separates object construction from
// formatting logic — renderers focus on Render/RenderMessage, while
// the factory owns construction policies (default options, icon sets, etc.).
//
// Usage:
//
//	factory := render.NewRendererFactory()
//	r := factory.Create("json")
//	factory.RegisterAll() // populates the global registry
type RendererFactory struct {
	// creators maps renderer names to their constructor functions.
	creators map[RendererName]func() Renderer
}

// NewRendererFactory returns a factory pre-loaded with all built-in
// renderer constructors. No renderers are registered globally until
// RegisterAll (or RegisterOne) is called.
func NewRendererFactory() *RendererFactory {
	f := &RendererFactory{
		creators: make(map[RendererName]func() Renderer),
	}
	f.registerBuiltins()
	return f
}

// registerBuiltins adds all built-in renderer constructors to the factory.
func (f *RendererFactory) registerBuiltins() {
	f.creators[RendererJSON] = func() Renderer { return NewJSONRenderer() }
	f.creators[RendererYAML] = func() Renderer { return NewYAMLRenderer() }
	f.creators[RendererPlain] = func() Renderer { return NewPlainRenderer() }
	f.creators[RendererColored] = func() Renderer { return NewColoredRenderer() }
	f.creators[RendererTable] = func() Renderer { return NewTableRenderer() }
	f.creators[RendererCompact] = func() Renderer { return NewCompactRenderer() }
}

// Create constructs a new renderer by name. Returns nil if the name
// is not known to the factory.
func (f *RendererFactory) Create(name RendererName) Renderer {
	creator, ok := f.creators[name]
	if !ok {
		return nil
	}
	return creator()
}

// RegisterAll creates and registers all known renderers in the global
// registry, replacing any previously registered renderer of the same name.
func (f *RendererFactory) RegisterAll() {
	for name, creator := range f.creators {
		_ = name // used implicitly via creator
		Register(creator())
	}
}

// RegisterOne creates a single renderer by name and registers it
// in the global registry. Returns false if the name is unknown.
func (f *RendererFactory) RegisterOne(name RendererName) bool {
	creator, ok := f.creators[name]
	if !ok {
		return false
	}
	Register(creator())
	return true
}

// AddCreator registers a custom renderer constructor with the factory.
// This allows third-party renderers to participate in the factory pattern.
func (f *RendererFactory) AddCreator(name RendererName, creator func() Renderer) {
	f.creators[name] = creator
}

// Names returns all renderer names known to this factory.
func (f *RendererFactory) Names() []RendererName {
	names := make([]RendererName, 0, len(f.creators))
	for name := range f.creators {
		names = append(names, name)
	}
	return names
}

// CreateColored constructs a ColoredRenderer with custom icons.
// This is a typed factory method for when callers need icon customization.
func (f *RendererFactory) CreateColored(icons Icons) *ColoredRenderer {
	return NewColoredRendererWithIcons(icons)
}

// init registers all built-in renderers via the factory.
// This replaces the per-file init() functions for a single registration point.
func init() {
	factory := NewRendererFactory()
	factory.RegisterAll()
}
