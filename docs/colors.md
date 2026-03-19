# colors

Import path: `github.com/rmkohlman/MaestroSDK/colors`

The `colors` package provides a decoupled interface for accessing theme colors in CLI tools. It defines the `ColorProvider` interface that consumers use instead of importing theme internals directly. This keeps the `render` package and other consumers decoupled from any specific theme implementation.

**Dependency flow:**

```
cmd/            -- injects ColorProvider via context
render/         -- uses ColorProvider interface (no theme import)
colors/         -- defines interface, implements via palette
MaestroPalette  -- pure data model (palette.Palette)
```

---

## ColorProvider Interface

```go
type ColorProvider interface {
    // Primary colors
    Primary() string    // Main brand/accent color
    Secondary() string  // Secondary brand color
    Accent() string     // Highlight/focus color

    // Status colors
    Success() string    // Green-ish success state
    Warning() string    // Yellow-ish warning state
    Error() string      // Red-ish error state
    Info() string       // Blue-ish info state

    // UI colors
    Foreground() string // Main text color
    Background() string // Main background color
    Muted() string      // Subdued/disabled text
    Highlight() string  // Selection/hover background
    Border() string     // Border/separator color

    // Theme metadata
    Name() string       // Theme name (e.g., "tokyonight-night")
    IsLight() bool      // Whether this is a light theme
}
```

All color methods return hex strings (e.g., `"#7aa2f7"`). When no color is applicable (as with `NoColorProvider`), methods return an empty string `""`.

---

## Implementations

### ThemeColorProvider

`ThemeColorProvider` implements `ColorProvider` using a `*palette.Palette` from the `MaestroPalette` module. It maps palette semantic color constants to `ColorProvider` methods and uses fallback chains when a palette key is absent.

```go
func NewThemeColorProvider(p *palette.Palette) ColorProvider
```

- If `p` is `nil`, returns a `NewDefaultColorProvider()` instead of panicking.
- Fallback chains try multiple palette keys in priority order (e.g., `Primary()` tries `palette.ColorPrimary` then `palette.ColorAccent` then the default hex value).

### NoColorProvider

`NoColorProvider` implements `ColorProvider` by returning empty strings for all color methods. Used when `--no-color` is passed or the `NO_COLOR` environment variable is set.

```go
type NoColorProvider struct{}

func NewNoColorProvider() ColorProvider
```

- `Name()` returns `"no-color"`.
- `IsLight()` returns `false`.
- All color methods return `""`.

### Default Providers

Two convenience constructors create `ColorProvider` instances backed by hardcoded default color maps (Tokyo Night inspired for dark, compatible defaults for light):

```go
func NewDefaultColorProvider() ColorProvider      // dark theme defaults
func NewDefaultLightColorProvider() ColorProvider // light theme defaults
```

**Default dark colors (`DefaultDarkColors`):**

| Key | Hex |
|-----|-----|
| primary | `#7aa2f7` |
| secondary | `#bb9af7` |
| accent | `#7aa2f7` |
| success | `#9ece6a` |
| warning | `#e0af68` |
| error | `#f7768e` |
| info | `#7dcfff` |
| foreground | `#c0caf5` |
| background | `#1a1b26` |
| muted | `#565f89` |
| highlight | `#283457` |
| border | `#414868` |

**Default light colors (`DefaultLightColors`):**

| Key | Hex |
|-----|-----|
| primary | `#3d59a1` |
| secondary | `#9854f1` |
| accent | `#3d59a1` |
| success | `#587539` |
| warning | `#8c6c3e` |
| error | `#c64343` |
| info | `#007197` |
| foreground | `#24292e` |
| background | `#ffffff` |
| muted | `#586069` |
| highlight | `#e1e4e8` |
| border | `#d0d7de` |

### MockColorProvider

`MockColorProvider` is a configurable implementation for use in tests. It is exported so test packages in other modules can use it directly.

```go
type MockColorProvider struct { /* unexported fields */ }

func NewMockColorProvider(opts ...MockOption) *MockColorProvider
```

Defaults to dark theme colors with name `"mock-theme"`.

**MockOption functional options:**

```go
func WithMockName(name string) MockOption
func WithMockLight() MockOption
func WithMockColor(key, value string) MockOption
func WithMockColors(colors map[string]string) MockOption
```

**Test helper methods on MockColorProvider:**

```go
func (m *MockColorProvider) SetColor(key, value string)
func (m *MockColorProvider) GetAllColors() map[string]string
```

`GetAllColors()` returns a copy of the internal map; mutations to the returned map do not affect the provider.

---

## Context Injection

The `colors` package provides functions for storing and retrieving a `ColorProvider` in a `context.Context`. This is the primary mechanism for passing colors through a command's call stack without explicit parameter threading.

```go
func WithProvider(ctx context.Context, provider ColorProvider) context.Context
```
Injects `provider` into `ctx`. Call this in the command's `PersistentPreRunE` or equivalent setup.

```go
func FromContext(ctx context.Context) (ColorProvider, bool)
```
Retrieves the provider. Returns `nil, false` if not present.

```go
func MustFromContext(ctx context.Context) ColorProvider
```
Retrieves the provider. Panics with `"colors: ColorProvider not found in context - did you call WithProvider?"` if not present.

```go
func FromContextOrDefault(ctx context.Context) ColorProvider
```
Returns the provider from context, or `NewDefaultColorProvider()` if not found.

```go
func FromContextOrDefaultLight(ctx context.Context) ColorProvider
```
Returns the provider from context, or `NewDefaultLightColorProvider()` if not found.

```go
func HasProvider(ctx context.Context) bool
```
Returns `true` if a `ColorProvider` is present in the context.

---

## Factory

### PaletteProvider Interface

`PaletteProvider` is an interface that bridges the factory to a theme system without a direct import. Callers implement this interface to supply palettes.

```go
type PaletteProvider interface {
    GetActivePalette() (*palette.Palette, error)
    GetPalette(name string) (*palette.Palette, error)
}
```

### ProviderFactory Interface

```go
type ProviderFactory interface {
    CreateFromActive() (ColorProvider, error)
    CreateFromTheme(themeName string) (ColorProvider, error)
    CreateDefault(isLight bool) ColorProvider
}

func NewProviderFactory(pp PaletteProvider) ProviderFactory
```

- `CreateFromActive()`: Loads the active theme from `pp`. On error, returns a usable `NewDefaultColorProvider()` and wraps the error as `"loading active theme: <err>"`.
- `CreateFromTheme(themeName)`: Returns an error if the theme cannot be found.
- `CreateDefault(isLight)`: Returns a dark or light default provider.

### Static Factory Functions

For simple use cases that do not need the full factory:

```go
func FromPalette(p *palette.Palette) ColorProvider
func Default() ColorProvider
func DefaultLight() ColorProvider
```

### NoProviderError

```go
type NoProviderError struct { /* unexported fields */ }
func (e *NoProviderError) Error() string
```

Returned by `CreateFromTheme` when no `PaletteProvider` is configured.

---

## CLI Helpers

These functions are the primary entry points for wiring up colors in CLI commands (e.g., cobra `PersistentPreRunE`).

### InitColorProviderForCommand

```go
func InitColorProviderForCommand(
    ctx context.Context,
    provider PaletteProvider,
    noColor bool,
) (context.Context, error)
```

Initializes a `ColorProvider` and injects it into the returned context. Behavior:

1. If `noColor` is `true` or `NO_COLOR` env var is set: injects `NoColorProvider`, returns `nil` error.
2. If `provider` is `nil`: injects `NewDefaultColorProvider()`, returns `nil` error.
3. Otherwise: creates a `ProviderFactory` and calls `CreateFromActive()`. The returned context always has a usable provider; the error is informational (theme load failure).

### InitColorProviderWithTheme

```go
func InitColorProviderWithTheme(
    ctx context.Context,
    provider PaletteProvider,
    themeName string,
    noColor bool,
) (context.Context, error)
```

Like `InitColorProviderForCommand` but loads a specific named theme. Returns an error (and the original context unchanged) if the theme cannot be loaded.

### GetDefaultThemePath

```go
func GetDefaultThemePath() string
```

Returns a theme storage path by checking in order:
1. `DVM_THEME_PATH` environment variable
2. `$XDG_CONFIG_HOME/dvm/themes`
3. `$HOME/.config/dvm/themes`
4. `./themes` (last resort)

### IsNoColorRequested

```go
func IsNoColorRequested(noColorFlag bool) bool
```

Returns `true` if `noColorFlag` is set, `NO_COLOR` env var is non-empty, or `TERM` equals `"dumb"`.

---

## PaletteAdapter

`PaletteAdapter` and `ColorToPaletteAdapter` implement the Adapter Pattern to convert a `ColorProvider` back to a `*palette.Palette` when palette-based rendering components (e.g., Starship prompt generators) need a full palette.

### PaletteAdapter Interface

```go
type PaletteAdapter interface {
    ToPalette() *palette.Palette
}
```

### ColorToPaletteAdapter

```go
type ColorToPaletteAdapter struct { /* unexported fields */ }

func NewColorToPaletteAdapter(provider ColorProvider) PaletteAdapter
```

Maps `ColorProvider` methods to palette semantic constants:

| ColorProvider method | Palette key |
|----------------------|-------------|
| `Primary()` | `palette.ColorPrimary` |
| `Secondary()` | `palette.ColorSecondary` |
| `Accent()` | `palette.ColorAccent` |
| `Success()` | `palette.ColorSuccess` |
| `Warning()` | `palette.ColorWarning` |
| `Error()` | `palette.ColorError` |
| `Info()` | `palette.ColorInfo` |
| `Foreground()` | `palette.ColorFg` |
| `Background()` | `palette.ColorBg` |
| `Muted()` | `palette.ColorComment` |
| `Highlight()` | `palette.ColorBgHighlight` |
| `Border()` | `palette.ColorBorder` |

Empty color values (as produced by `NoColorProvider`) are omitted from the resulting palette rather than stored as empty strings.

### ToPalette Convenience Function

```go
func ToPalette(provider ColorProvider) *palette.Palette
```

Creates a `ColorToPaletteAdapter` and calls `ToPalette()` in one call. Equivalent to:

```go
adapter := colors.NewColorToPaletteAdapter(provider)
p := adapter.ToPalette()
```
