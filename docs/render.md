# render

The `render` package provides a decoupled rendering system for CLI output. It separates data preparation from display logic: commands prepare structured data and pass it to the renderer, which decides how to display it based on the active renderer and output options.

**Architecture:**

```
Command layer (prepares data)
        |
        v
render package (Output / OutputWithContext)
        |
        v
Renderer interface (JSON, YAML, Colored, Plain, Table, Compact)
        |
        v
io.Writer (stdout, file, buffer)
```

**Renderer selection priority (highest to lowest):**

1. Explicit override passed to `OutputWithContextAndRenderer` / `OutputWith`
2. `DVM_RENDER` environment variable
3. Global config default (set via `SetDefault` or `SetConfig`)
4. Built-in default: `colored`

---

## Types

### RenderType

`RenderType` hints to the renderer what kind of data structure is being passed.

| Constant | Value | Description |
|----------|-------|-------------|
| `TypeAuto` | `"auto"` | Renderer infers the best display format |
| `TypeKeyValue` | `"keyvalue"` | Key-value pairs |
| `TypeTable` | `"table"` | Tabular data with headers and rows |
| `TypeList` | `"list"` | Simple list of strings |
| `TypeDetail` | `"detail"` | Single-item detail view |
| `TypeRaw` | `"raw"` | Pass-through without formatting |
| `TypeProgress` | `"progress"` | Progress indicator output |

### RendererName

| Constant | Value |
|----------|-------|
| `RendererJSON` | `"json"` |
| `RendererYAML` | `"yaml"` |
| `RendererColored` | `"colored"` |
| `RendererPlain` | `"plain"` |
| `RendererTable` | `"table"` |
| `RendererCompact` | `"compact"` |

### Options

`Options` configures how data should be rendered. Commands set these to provide hints to the renderer.

| Field | Type | Description |
|-------|------|-------------|
| `Type` | `RenderType` | Data structure hint |
| `Title` | `string` | Section title (human-readable renderers only) |
| `Headers` | `[]string` | Column headers for table type |
| `Empty` | `bool` | Indicates the data represents an empty state |
| `EmptyMessage` | `string` | Message shown when `Empty` is true |
| `EmptyHints` | `[]string` | Suggestions shown when `Empty` is true |
| `Verbose` | `bool` | Enable extra detail |
| `Wide` | `bool` | Enable wide format with additional columns |

### MessageLevel

| Constant | Value |
|----------|-------|
| `LevelInfo` | `"info"` |
| `LevelSuccess` | `"success"` |
| `LevelWarning` | `"warning"` |
| `LevelError` | `"error"` |
| `LevelDebug` | `"debug"` |
| `LevelProgress` | `"progress"` |

### Message

A `Message` combines a `MessageLevel` and a content string. Pass to `RenderMessage` variants.

### Config

`Config` holds global renderer configuration: the default renderer name, verbose flag, no-color flag, and Nerd Font icon flag. `DefaultConfig()` returns a config with `RendererColored` as the default.

### Data Structures

Commands pass these types as the `data` argument to `Output`/`OutputWithContext`. Renderers handle each type appropriately.

| Type | Fields | Description |
|------|--------|-------------|
| `TableData` | `Headers []string`, `Rows [][]string` | Pre-formatted table |
| `ListData` | `Items []string` | Simple string list |
| `KeyValueData` | `Pairs []KeyValue` | Ordered key-value pairs |
| `KeyValue` | `Key string`, `Value string` | Single key-value entry |

`NewKeyValueData(map)` builds `KeyValueData` from a map (key order not guaranteed). `NewOrderedKeyValueData(pairs...)` preserves explicit ordering.

---

## Renderer Interface

All six renderers implement a common interface exposing `Render`, `RenderWithContext`, `RenderMessage`, `RenderMessageWithContext`, `Name`, and `SupportsColor`. The non-context variants call their `WithContext` equivalents using `context.Background()` for backward compatibility.

---

## Renderers

### ColoredRenderer

The default renderer. Outputs richly formatted text with lipgloss styles and Unicode icons. When a `ColorProvider` is available in the context, it derives its styles from that provider; otherwise it falls back to built-in Catppuccin-inspired defaults.

- `Name()` returns `RendererColored`
- `SupportsColor()` returns `true`
- Registered automatically at startup

Handles: `KeyValueData`, `TableData`, `ListData`, `[]string`, `map[string]string`, `map[string]interface{}`, and any other type via `fmt.Fprintf`.

Three icon sets are provided: `DefaultIcons()` (Unicode), `NerdFontIcons()` (Nerd Font codepoints), and `PlainIcons()` (ASCII).

### PlainRenderer

Outputs plain text without color or styling. Suitable for piping, CI environments, and terminals without color support. Ignores context entirely.

- `Name()` returns `RendererPlain`
- `SupportsColor()` returns `false`
- Message prefixes: `[OK]`, `[WARN]`, `[ERROR]`, `[DEBUG]`, `->` (progress), `[INFO]`

### JSONRenderer

Outputs data as indented JSON (2-space indent). Ignores `Title`, `EmptyMessage`, `EmptyHints`, and context. Converts `KeyValueData` to a key-value map, `TableData` to a list of row maps, and `ListData` to a string array. Messages are output as `{"level": "...", "message": "..."}`.

- `Name()` returns `RendererJSON`
- `SupportsColor()` returns `false`

### YAMLRenderer

Outputs data as YAML (2-space indent). Behavior mirrors `JSONRenderer` for data type conversions.

- `Name()` returns `RendererYAML`
- `SupportsColor()` returns `false`

### TableRenderer

Focuses exclusively on table output. Suppresses all `RenderMessage` calls (no-op). Skips output entirely when `opts.Empty` is `true`. Delegates actual table rendering to an embedded `ColoredRenderer`, so it does respond to context-provided colors.

- `Name()` returns `RendererTable`
- `SupportsColor()` returns `true`

For `KeyValueData` and `map[string]string`, renders as `key: value` lines. For non-table data types, outputs nothing.

### CompactRenderer

Like `ColoredRenderer` but more condensed: tighter column spacing, muted-style headers (no separator line), `▸` title prefix instead of `▌`, and compact list items use `-` bullet instead of `•`. Delegates to `ColoredRenderer` for unsupported data types and for message rendering.

- `Name()` returns `RendererCompact`
- `SupportsColor()` returns `true`

---

## Registry

The global registry maps `RendererName` to `Renderer` implementations. All six built-in renderers register themselves at startup.

**Registration:** `Register(r Renderer)` — adds a renderer. `Get(name)` — retrieves by name. `List()` — returns all registered names.

**Configuration:** `SetConfig(cfg Config)`, `GetConfig()`, `SetDefault(name RendererName)`.

**Writer:** `SetWriter(w io.Writer)`, `GetWriter() io.Writer`. The default writer is stdout.

**Renderer resolution:** `ResolveRenderer(override string)` applies the priority chain: override string → `DVM_RENDER` env var → global config default. If `NO_COLOR` is set and the resolved renderer is `colored` or `compact`, it falls back to `plain`. If the resolved renderer is not registered, falls back to `colored`, then `plain`.

---

## Output Functions

### Context-Aware (Preferred)

| Function | Description |
|----------|-------------|
| `OutputWithContext` | Renders data using the resolved renderer; passes context for color access |
| `OutputWithContextAndRenderer` | Renders data with an explicit renderer override |
| `OutputToWithContext` | Renders data to a specific writer with a renderer override |

### Non-Context (Backward Compatible)

| Function | Description |
|----------|-------------|
| `Output` | Renders data using the resolved renderer |
| `OutputWith` | Renders data with a renderer override |
| `OutputTo` | Renders data to a specific writer |

---

## Message Functions

### Context-Aware

| Function | Description |
|----------|-------------|
| `MsgWithContext` | Renders a message at the given level using the context-aware renderer |
| `MsgWithContextAndRenderer` | Renders a message with a renderer override |
| `MsgToWithContext` | Renders a message to a specific writer |

### Non-Context

| Function | Description |
|----------|-------------|
| `Msg` | Renders a message at the given level |
| `MsgWith` | Renders a message with a renderer override |
| `MsgTo` | Renders a message to a specific writer |

### Convenience Message Functions

`Info`, `Success`, `Warning`, `Error`, and `Progress` call `Msg` with the corresponding level. Formatted variants (`Infof`, `Successf`, `Warningf`, `Errorf`, `Progressf`) accept a format string and args.

Stderr output variants: `InfoToStderr`, `WarningToStderr`, `ErrorToStderr`, and their `f` counterparts.

Undecorated text output: `Plain(text)` and `Plainf(format, args...)` — no level prefix, no color. `Blank()` outputs an empty line.
