# colors

The `colors` package provides a decoupled interface for accessing theme colors in CLI tools. It defines the `ColorProvider` interface that consumers use instead of importing theme internals directly. This keeps the `render` package and other consumers decoupled from any specific theme implementation.

**Dependency flow:**

```
CLI commands      -- inject ColorProvider via context
render package    -- uses ColorProvider interface (no theme import)
colors package    -- defines interface, implements via palette
MaestroPalette    -- pure data model (palette.Palette)
```

---

## ColorProvider Interface

`ColorProvider` is the primary interface. All color methods return hex strings (e.g., `"#7aa2f7"`). When no color is applicable (as with `NoColorProvider`), methods return an empty string `""`.

| Method | Description |
|--------|-------------|
| `Primary()` | Main brand/accent color |
| `Secondary()` | Secondary brand color |
| `Accent()` | Highlight/focus color |
| `Success()` | Green-ish success state |
| `Warning()` | Yellow-ish warning state |
| `Error()` | Red-ish error state |
| `Info()` | Blue-ish info state |
| `Foreground()` | Main text color |
| `Background()` | Main background color |
| `Muted()` | Subdued/disabled text |
| `Highlight()` | Selection/hover background |
| `Border()` | Border/separator color |
| `Name()` | Theme name (e.g., `"tokyonight-night"`) |
| `IsLight()` | Whether this is a light theme |

---

## Implementations

### ThemeColorProvider

`ThemeColorProvider` implements `ColorProvider` using a palette from the `MaestroPalette` module. It maps palette semantic color constants to `ColorProvider` methods and uses fallback chains when a palette key is absent. If the palette is `nil`, returns a default color provider instead of panicking.

### NoColorProvider

`NoColorProvider` implements `ColorProvider` by returning empty strings for all color methods. Used when `--no-color` is passed or the `NO_COLOR` environment variable is set.

- `Name()` returns `"no-color"`
- `IsLight()` returns `false`
- All color methods return `""`

### Default Providers

Two convenience constructors create `ColorProvider` instances backed by hardcoded default color maps (Tokyo Night-inspired for dark, compatible defaults for light):

- `NewDefaultColorProvider()` — dark theme defaults
- `NewDefaultLightColorProvider()` — light theme defaults

**Default dark colors:**

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

**Default light colors:**

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

`MockColorProvider` is a configurable implementation for use in tests. It is exported so test packages in other modules can use it directly. Defaults to dark theme colors with name `"mock-theme"`.

Available options when constructing: `WithMockName`, `WithMockLight`, `WithMockColor`, `WithMockColors`. Test helper methods: `SetColor(key, value)`, `GetAllColors()`.

---

## Context Injection

The `colors` package provides functions for storing and retrieving a `ColorProvider` in a `context.Context`. This is the primary mechanism for passing colors through a command's call stack without explicit parameter threading.

| Function | Description |
|----------|-------------|
| `WithProvider(ctx, provider)` | Injects `provider` into `ctx`. Call this in the command's setup. |
| `FromContext(ctx)` | Retrieves the provider. Returns `nil, false` if not present. |
| `MustFromContext(ctx)` | Retrieves the provider. Panics if not present. |
| `FromContextOrDefault(ctx)` | Returns the provider from context, or `NewDefaultColorProvider()` if not found. |
| `FromContextOrDefaultLight(ctx)` | Returns the provider from context, or the light default if not found. |
| `HasProvider(ctx)` | Returns `true` if a `ColorProvider` is present in the context. |

---

## Factory

### PaletteProvider Interface

`PaletteProvider` bridges the factory to a theme system without a direct import. Callers implement this interface to supply palettes. It exposes `GetActivePalette()` and `GetPalette(name)`.

### ProviderFactory

`ProviderFactory` creates `ColorProvider` instances from a `PaletteProvider`. Construct with `NewProviderFactory(pp PaletteProvider)`.

| Method | Description |
|--------|-------------|
| `CreateFromActive()` | Loads the active theme from the palette provider. On error, returns a usable default provider. |
| `CreateFromTheme(themeName)` | Returns an error if the theme cannot be found. |
| `CreateDefault(isLight)` | Returns a dark or light default provider. |

### Static Factory Functions

For simple use cases that do not need the full factory:

| Function | Description |
|----------|-------------|
| `FromPalette(p)` | Creates a `ColorProvider` from a palette |
| `Default()` | Returns the default dark color provider |
| `DefaultLight()` | Returns the default light color provider |

---

## CLI Helpers

These functions are the primary entry points for wiring up colors in CLI commands.

### InitColorProviderForCommand

Initializes a `ColorProvider` and injects it into the returned context.

1. If `noColor` is `true` or `NO_COLOR` env var is set: injects `NoColorProvider`, returns `nil` error.
2. If `provider` is `nil`: injects `NewDefaultColorProvider()`, returns `nil` error.
3. Otherwise: creates a `ProviderFactory` and loads from the active theme. The returned context always has a usable provider; any error is informational.

### InitColorProviderWithTheme

Like `InitColorProviderForCommand` but loads a specific named theme. Returns an error (and the original context unchanged) if the theme cannot be loaded.

### GetDefaultThemePath

Returns a theme storage path by checking in order:
1. `DVM_THEME_PATH` environment variable
2. `$XDG_CONFIG_HOME/dvm/themes`
3. `$HOME/.config/dvm/themes`
4. `./themes` (last resort)

### IsNoColorRequested

Returns `true` if the no-color flag is set, `NO_COLOR` env var is non-empty, or `TERM` equals `"dumb"`.

---

## PaletteAdapter

`PaletteAdapter` and `ColorToPaletteAdapter` implement the Adapter Pattern to convert a `ColorProvider` back to a palette when palette-based rendering components (e.g., Starship prompt generators) need a full palette.

`NewColorToPaletteAdapter(provider)` creates an adapter. Call `ToPalette()` on it to produce the palette.

`ToPalette(provider)` is a convenience function that creates the adapter and calls `ToPalette()` in one step.

| ColorProvider method | Palette key |
|----------------------|-------------|
| `Primary()` | `ColorPrimary` |
| `Secondary()` | `ColorSecondary` |
| `Accent()` | `ColorAccent` |
| `Success()` | `ColorSuccess` |
| `Warning()` | `ColorWarning` |
| `Error()` | `ColorError` |
| `Info()` | `ColorInfo` |
| `Foreground()` | `ColorFg` |
| `Background()` | `ColorBg` |
| `Muted()` | `ColorComment` |
| `Highlight()` | `ColorBgHighlight` |
| `Border()` | `ColorBorder` |

Empty color values (as produced by `NoColorProvider`) are omitted from the resulting palette rather than stored as empty strings.
