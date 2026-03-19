# API Reference

Quick-lookup for all exported symbols in MaestroSDK. Each entry links to its detailed documentation page.

---

## colors

Import: `github.com/rmkohlman/MaestroSDK/colors`

Full documentation: [colors](colors.md)

### Interface

| Symbol | Kind | Description |
|--------|------|-------------|
| `ColorProvider` | interface | Provides hex color strings for all theme roles |

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `ThemeColorProvider` | struct | Wraps a `palette.Palette`; returned by `NewThemeColorProvider` |
| `NoColorProvider` | struct | Returns empty strings for all color methods; used when `NO_COLOR` is set |
| `PaletteAdapter` | struct | Adapts an external palette-like struct to `ColorProvider` via reflection |
| `NoProviderError` | struct | Error type returned when no `ColorProvider` is found in context |

### Functions

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `NewThemeColorProvider` | `(p *palette.Palette) *ThemeColorProvider` | Creates a `ColorProvider` backed by a `MaestroPalette` palette |
| `NewNoColorProvider` | `() *NoColorProvider` | Creates a no-op `ColorProvider` |
| `NewPaletteAdapter` | `(src any) (*PaletteAdapter, error)` | Creates a `PaletteAdapter` from any struct with matching color fields |
| `DefaultProvider` | `() ColorProvider` | Returns a built-in default `ColorProvider` (suitable as a fallback) |
| `DefaultNoColorProvider` | `() *NoColorProvider` | Returns a package-level `NoColorProvider` singleton |
| `WithColorProvider` | `(ctx context.Context, p ColorProvider) context.Context` | Injects `ColorProvider` into a `context.Context` |
| `FromContext` | `(ctx context.Context) (ColorProvider, error)` | Retrieves `ColorProvider` from context; returns `NoProviderError` if absent |
| `FromContextOrDefault` | `(ctx context.Context) ColorProvider` | Retrieves `ColorProvider` from context, falling back to `DefaultProvider()` |
| `InitColorProviderForCommand` | `(ctx context.Context, pp PaletteProvider, noColor bool) (context.Context, error)` | Resolves and injects `ColorProvider` for a CLI command |
| `NewMockProvider` | `() *MockProvider` | Creates a `MockProvider` pre-filled with recognizable test colors |

### Interfaces (supporting)

| Symbol | Kind | Description |
|--------|------|-------------|
| `PaletteProvider` | interface | Single method `GetPalette() (*palette.Palette, error)`; implemented by palette loaders |

---

## render

Import: `github.com/rmkohlman/MaestroSDK/render`

Full documentation: [render](render.md)

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `RenderType` | `string` type | Hints the renderer about data structure (`auto`, `keyvalue`, `table`, `list`, `detail`, `raw`, `progress`) |
| `RendererName` | `string` type | Identifies a renderer (`json`, `yaml`, `colored`, `plain`, `table`, `compact`) |
| `Options` | struct | Rendering configuration: `Type`, `Title`, `Headers`, `Empty`, `EmptyMessage`, `EmptyHints`, `Verbose`, `Wide` |
| `Config` | struct | Global renderer configuration: `Default RendererName`, `Writer io.Writer` |
| `TableData` | struct | Pre-formatted table: `Headers []string`, `Rows [][]string` |
| `Renderer` | interface | Core interface: `Render`, `RenderWithContext`, `RenderMessage`, `RenderMessageWithContext`, `Name` |

### RenderType Constants

| Constant | Value |
|----------|-------|
| `TypeAuto` | `"auto"` |
| `TypeKeyValue` | `"keyvalue"` |
| `TypeTable` | `"table"` |
| `TypeList` | `"list"` |
| `TypeDetail` | `"detail"` |
| `TypeRaw` | `"raw"` |
| `TypeProgress` | `"progress"` |

### RendererName Constants

| Constant | Value |
|----------|-------|
| `RendererJSON` | `"json"` |
| `RendererYAML` | `"yaml"` |
| `RendererColored` | `"colored"` |
| `RendererPlain` | `"plain"` |
| `RendererTable` | `"table"` |
| `RendererCompact` | `"compact"` |

### Renderer Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `ColoredRenderer` | struct | Human-readable output with lipgloss styling; default renderer |
| `PlainRenderer` | struct | Human-readable output with no color or styling |
| `JSONRenderer` | struct | Marshals data to JSON |
| `YAMLRenderer` | struct | Marshals data to YAML |
| `TableRenderer` | struct | Tabular output; `RenderMessage` is a deliberate no-op |
| `CompactRenderer` | struct | Embeds `*ColoredRenderer`; emits compact one-liner output |

### Output Functions (primary API)

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `Output` | `(data any, opts Options) error` | Renders data using the resolved renderer |
| `OutputWithContext` | `(ctx context.Context, data any, opts Options) error` | Renders data; passes context to renderer for `ColorProvider` access |
| `OutputWithContextAndRenderer` | `(ctx context.Context, data any, opts Options, rendererName string) error` | Renders data with an explicit renderer override |
| `OutputWith` | `(ctx context.Context, data any, opts Options, rendererName string, w io.Writer) error` | Renders data to a specific `io.Writer` with a renderer override |
| `Message` | `(msg string) error` | Renders a plain message string using the resolved renderer |
| `MessageWithContext` | `(ctx context.Context, msg string) error` | Renders a plain message string using the context-aware renderer |

### Registry Functions

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `Register` | `(r Renderer)` | Registers a renderer; called from each renderer's `init()` |
| `ResolveRenderer` | `(name string) Renderer` | Returns the named renderer, applying `NO_COLOR` fallback logic |
| `SetDefault` | `(name RendererName)` | Sets the global default renderer |
| `SetConfig` | `(cfg Config)` | Sets global renderer config (default name + writer) |
| `GetWriter` | `() io.Writer` | Returns the active output writer (defaults to `os.Stdout`) |

### Convenience Constructors

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `NewColoredRenderer` | `(w io.Writer) *ColoredRenderer` | Creates a `ColoredRenderer` writing to `w` |
| `NewPlainRenderer` | `(w io.Writer) *PlainRenderer` | Creates a `PlainRenderer` writing to `w` |
| `NewJSONRenderer` | `(w io.Writer) *JSONRenderer` | Creates a `JSONRenderer` writing to `w` |
| `NewYAMLRenderer` | `(w io.Writer) *YAMLRenderer` | Creates a `YAMLRenderer` writing to `w` |
| `NewTableRenderer` | `(w io.Writer) *TableRenderer` | Creates a `TableRenderer` writing to `w` |
| `NewCompactRenderer` | `(w io.Writer) *CompactRenderer` | Creates a `CompactRenderer` writing to `w` |

---

## resource

Import: `github.com/rmkohlman/MaestroSDK/resource`

Full documentation: [resource](resource.md)

### Interfaces

| Symbol | Kind | Description |
|--------|------|-------------|
| `Resource` | interface | `GetKind() string`, `GetName() string`, `Validate() error` |
| `Handler` | interface | CRUD for a single resource kind: `Kind`, `Apply`, `Get`, `List`, `Delete`, `ToYAML` |

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `Context` | struct | Dependency container for handlers: `DataStore any`, `PluginStore any`, `ThemeStore any`, `ConfigDir string` |
| `KindHeader` | struct | Used internally to detect `Kind` from YAML before full parsing |
| `ResourceList` | struct | Ordered collection of `Resource` values with dependency-aware sorting |

### Registry Functions

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `Register` | `(h Handler)` | Registers a handler for its `Kind()` |
| `Apply` | `(ctx Context, data []byte, source string) (Resource, error)` | Parses YAML, resolves handler, calls `Apply` |
| `Get` | `(ctx Context, kind, name string) (Resource, error)` | Retrieves one resource by kind and name |
| `List` | `(ctx Context, kind string) ([]Resource, error)` | Lists all resources of a kind |
| `Delete` | `(ctx Context, kind, name string) error` | Deletes a resource by kind and name |
| `ToYAML` | `(ctx Context, kind, name string) ([]byte, error)` | Serializes a resource to YAML |
| `GetHandler` | `(kind string) (Handler, error)` | Returns the registered handler for a kind |
| `RegisteredKinds` | `() []string` | Returns all registered kind names, sorted |

### ResourceList Functions

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `NewResourceList` | `(items []Resource) *ResourceList` | Creates a `ResourceList` from a slice |
| `(rl *ResourceList) DependencyOrder` | `(depFn func(Resource) []string) ([]Resource, error)` | Topologically sorts resources by dependency |
| `(rl *ResourceList) Items` | `() []Resource` | Returns the underlying slice |

### Generic DataStore Helpers

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `DataStoreAs[T]` | `(ctx Context) (T, error)` | Type-asserts `ctx.DataStore` to `T` |
| `PluginStoreAs[T]` | `(ctx Context) (T, error)` | Type-asserts `ctx.PluginStore` to `T` |
| `ThemeStoreAs[T]` | `(ctx Context) (T, error)` | Type-asserts `ctx.ThemeStore` to `T` |

---

## paths

Import: `github.com/rmkohlman/MaestroSDK/paths`

Full documentation: [paths](paths.md)

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `DVMDirName` | `".devopsmaestro"` | Hidden directory under `$HOME` for `dvm` state |
| `NVPDirName` | `".nvp"` | Hidden directory under `$HOME` for `nvp` state |
| `DVTDirName` | `".dvt"` | Hidden directory under `$HOME` for `dvt` state |
| `DatabaseFile` | `"devopsmaestro.db"` | SQLite database filename inside the `dvm` root |

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `PathConfig` | struct | Immutable path resolver rooted at a home directory |

### Constructors

| Symbol | Signature | Description |
|--------|-----------|-------------|
| `New` | `(homeDir string) *PathConfig` | Creates a `PathConfig` for `homeDir`; panics if empty |
| `Default` | `() (*PathConfig, error)` | Creates a `PathConfig` using `os.UserHomeDir()` |

### PathConfig Methods — DVM Root

| Method | Returns |
|--------|---------|
| `Root()` | `{home}/.devopsmaestro` |
| `ConfigFile()` | `{root}/config.yaml` |
| `Database()` | `{root}/devopsmaestro.db` |
| `VersionFile()` | `{root}/.version` |
| `ContextFile()` | `{root}/context.yaml` |
| `NvimSyncStatus()` | `{root}/.nvim-sync-status` |
| `LogsDir()` | `{root}/logs` |
| `BackupsDir()` | `{root}/backups` |
| `TemplatesDir()` | `{root}/templates` |
| `NvimTemplatesDir()` | `{root}/templates/nvim` |
| `ShellTemplatesDir()` | `{root}/templates/shell` |

### PathConfig Methods — Workspaces

| Method | Signature | Returns |
|--------|-----------|---------|
| `WorkspacesDir()` | `() string` | `{root}/workspaces` |
| `WorkspacePath()` | `(slug string) string` | `{root}/workspaces/{slug}` |
| `WorkspaceRepoPath()` | `(slug string) string` | `{root}/workspaces/{slug}/repo` |
| `WorkspaceVolumePath()` | `(slug string) string` | `{root}/workspaces/{slug}/volume` |
| `WorkspaceConfigPath()` | `(slug string) string` | `{root}/workspaces/{slug}/.dvm` |

### PathConfig Methods — Git and Build

| Method | Signature | Returns |
|--------|-----------|---------|
| `ReposDir()` | `() string` | `{root}/repos` |
| `BuildStagingDir()` | `(appName string) string` | `{root}/build-staging/{appName}` |

### PathConfig Methods — Registry

| Method | Signature | Returns |
|--------|-----------|---------|
| `RegistryDir()` | `(name string) string` | `{root}/registries/{name}` |
| `RegistryStorage()` | `() string` | `{root}/registry` |
| `AthensStorage()` | `() string` | `{root}/athens` |
| `VerdaccioStorage()` | `() string` | `{root}/verdaccio` |
| `DevpiStorage()` | `() string` | `{root}/devpi` |
| `SquidDir()` | `() string` | `{root}/squid` |

### PathConfig Methods — NVP

| Method | Returns |
|--------|---------|
| `NVPRoot()` | `{home}/.nvp` |
| `NVPPluginsDir()` | `{nvpRoot}/plugins` |
| `NVPPackagesDir()` | `{nvpRoot}/packages` |
| `NVPThemesDir()` | `{nvpRoot}/themes` |
| `NVPCoreConfig()` | `{nvpRoot}/core.yaml` |

### PathConfig Methods — DVT

| Method | Returns |
|--------|---------|
| `DVTRoot()` | `{home}/.dvt` |
| `DVTPromptsDir()` | `{dvtRoot}/prompts` |
| `DVTPluginsDir()` | `{dvtRoot}/plugins` |
| `DVTShellsDir()` | `{dvtRoot}/shells` |
| `DVTProfilesDir()` | `{dvtRoot}/profiles` |
| `DVTActiveProfile()` | `{dvtRoot}/.active-profile` |

### PathConfig Methods — Helpers

| Method | Returns | Notes |
|--------|---------|-------|
| `DatabasePathTilde()` | `~/.devopsmaestro/devopsmaestro.db` | Tilde notation only — not a real filesystem path |
