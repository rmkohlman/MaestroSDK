# API Reference

Quick-lookup for all exported symbols in MaestroSDK. Each entry links to its detailed documentation page.

---

## colors

Full documentation: [colors](colors.md)

### Interface

| Symbol | Kind | Description |
|--------|------|-------------|
| `ColorProvider` | interface | Provides hex color strings for all theme roles |

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `ThemeColorProvider` | struct | Wraps a palette; returned by `NewThemeColorProvider` |
| `NoColorProvider` | struct | Returns empty strings for all color methods; used when `NO_COLOR` is set |
| `PaletteAdapter` | struct | Adapts an external palette-like struct to `ColorProvider` |
| `NoProviderError` | struct | Error type returned when no `ColorProvider` is found in context |

### Functions

| Symbol | Description |
|--------|-------------|
| `NewThemeColorProvider` | Creates a `ColorProvider` backed by a `MaestroPalette` palette |
| `NewNoColorProvider` | Creates a no-op `ColorProvider` |
| `NewPaletteAdapter` | Creates a `PaletteAdapter` from any struct with matching color fields |
| `DefaultProvider` | Returns a built-in default `ColorProvider` (suitable as a fallback) |
| `DefaultNoColorProvider` | Returns a package-level `NoColorProvider` singleton |
| `WithColorProvider` | Injects `ColorProvider` into a `context.Context` |
| `FromContext` | Retrieves `ColorProvider` from context; returns `NoProviderError` if absent |
| `FromContextOrDefault` | Retrieves `ColorProvider` from context, falling back to `DefaultProvider()` |
| `InitColorProviderForCommand` | Resolves and injects `ColorProvider` for a CLI command |
| `NewMockProvider` | Creates a `MockProvider` pre-filled with recognizable test colors |

### Supporting Interfaces

| Symbol | Kind | Description |
|--------|------|-------------|
| `PaletteProvider` | interface | Single method `GetPalette()`; implemented by palette loaders |

---

## render

Full documentation: [render](render.md)

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `RenderType` | string type | Hints the renderer about data structure (`auto`, `keyvalue`, `table`, `list`, `detail`, `raw`, `progress`) |
| `RendererName` | string type | Identifies a renderer (`json`, `yaml`, `colored`, `plain`, `table`, `compact`) |
| `Options` | struct | Rendering configuration: `Type`, `Title`, `Headers`, `Empty`, `EmptyMessage`, `EmptyHints`, `Verbose`, `Wide` |
| `Config` | struct | Global renderer configuration: `Default RendererName`, `Writer` |
| `TableData` | struct | Pre-formatted table: `Headers`, `Rows` |
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
| `CompactRenderer` | struct | Compact variant of `ColoredRenderer` |

### Output Functions (primary API)

| Symbol | Description |
|--------|-------------|
| `Output` | Renders data using the resolved renderer |
| `OutputWithContext` | Renders data; passes context to renderer for `ColorProvider` access |
| `OutputWithContextAndRenderer` | Renders data with an explicit renderer override |
| `OutputWith` | Renders data to a specific writer with a renderer override |
| `Message` | Renders a plain message string using the resolved renderer |
| `MessageWithContext` | Renders a plain message string using the context-aware renderer |

### Registry Functions

| Symbol | Description |
|--------|-------------|
| `Register` | Registers a renderer |
| `ResolveRenderer` | Returns the named renderer, applying `NO_COLOR` fallback logic |
| `SetDefault` | Sets the global default renderer |
| `SetConfig` | Sets global renderer config |
| `GetWriter` | Returns the active output writer (defaults to stdout) |

### Convenience Constructors

| Symbol | Description |
|--------|-------------|
| `NewColoredRenderer` | Creates a `ColoredRenderer` |
| `NewPlainRenderer` | Creates a `PlainRenderer` |
| `NewJSONRenderer` | Creates a `JSONRenderer` |
| `NewYAMLRenderer` | Creates a `YAMLRenderer` |
| `NewTableRenderer` | Creates a `TableRenderer` |
| `NewCompactRenderer` | Creates a `CompactRenderer` |

---

## resource

Full documentation: [resource](resource.md)

### Interfaces

| Symbol | Kind | Description |
|--------|------|-------------|
| `Resource` | interface | `GetKind()`, `GetName()`, `Validate()` |
| `Handler` | interface | CRUD for a single resource kind: `Kind`, `Apply`, `Get`, `List`, `Delete`, `ToYAML` |

### Types

| Symbol | Kind | Description |
|--------|------|-------------|
| `Context` | struct | Dependency container for handlers: `DataStore`, `PluginStore`, `ThemeStore`, `ConfigDir` |
| `KindHeader` | struct | Used internally to detect `Kind` from YAML before full parsing |
| `ResourceList` | struct | Ordered collection of `Resource` values with dependency-aware sorting |

### Registry Functions

| Symbol | Description |
|--------|-------------|
| `Register` | Registers a handler for its `Kind()` |
| `Apply` | Parses YAML, resolves handler, calls `Apply` |
| `Get` | Retrieves one resource by kind and name |
| `List` | Lists all resources of a kind |
| `Delete` | Deletes a resource by kind and name |
| `ToYAML` | Serializes a resource to YAML |
| `GetHandler` | Returns the registered handler for a kind |
| `RegisteredKinds` | Returns all registered kind names, sorted |

### ResourceList Functions

| Symbol | Description |
|--------|-------------|
| `NewResourceList` | Creates a `ResourceList` from a slice |
| `DependencyOrder` | Topologically sorts resources by dependency |
| `Items` | Returns the underlying slice |

### Generic DataStore Helpers

| Symbol | Description |
|--------|-------------|
| `DataStoreAs` | Type-asserts `ctx.DataStore` to a given type |
| `PluginStoreAs` | Type-asserts `ctx.PluginStore` to a given type |
| `ThemeStoreAs` | Type-asserts `ctx.ThemeStore` to a given type |

---

## paths

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

| Symbol | Description |
|--------|-------------|
| `New` | Creates a `PathConfig` for a given home directory; panics if empty |
| `Default` | Creates a `PathConfig` using `os.UserHomeDir()` |

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

| Method | Returns |
|--------|---------|
| `WorkspacesDir()` | `{root}/workspaces` |
| `WorkspacePath(slug)` | `{root}/workspaces/{slug}` |
| `WorkspaceRepoPath(slug)` | `{root}/workspaces/{slug}/repo` |
| `WorkspaceVolumePath(slug)` | `{root}/workspaces/{slug}/volume` |
| `WorkspaceConfigPath(slug)` | `{root}/workspaces/{slug}/.dvm` |

### PathConfig Methods — Git and Build

| Method | Returns |
|--------|---------|
| `ReposDir()` | `{root}/repos` |
| `BuildStagingDir(appName)` | `{root}/build-staging/{appName}` |

### PathConfig Methods — Registry

| Method | Returns |
|--------|---------|
| `RegistryDir(name)` | `{root}/registries/{name}` |
| `RegistryStorage()` | `{root}/registry` |
| `AthensStorage()` | `{root}/athens` |
| `VerdaccioStorage()` | `{root}/verdaccio` |
| `DevpiStorage()` | `{root}/devpi` |
| `SquidDir()` | `{root}/squid` |

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
