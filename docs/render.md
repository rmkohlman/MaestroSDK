# render

Import path: `github.com/rmkohlman/MaestroSDK/render`

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

```go
type RenderType string

const (
    TypeAuto     RenderType = "auto"
    TypeKeyValue RenderType = "keyvalue"
    TypeTable    RenderType = "table"
    TypeList     RenderType = "list"
    TypeDetail   RenderType = "detail"
    TypeRaw      RenderType = "raw"
    TypeProgress RenderType = "progress"
)
```

### RendererName

```go
type RendererName string

const (
    RendererJSON    RendererName = "json"
    RendererYAML    RendererName = "yaml"
    RendererColored RendererName = "colored"
    RendererPlain   RendererName = "plain"
    RendererTable   RendererName = "table"
    RendererCompact RendererName = "compact"
)
```

### Options

`Options` configures how data should be rendered. Commands set these to provide hints to the renderer.

```go
type Options struct {
    Type         RenderType // Data structure hint
    Title        string     // Section title (human-readable renderers only)
    Headers      []string   // Column headers for table type
    Empty        bool       // Indicates the data represents an empty state
    EmptyMessage string     // Message shown when Empty is true
    EmptyHints   []string   // Suggestions shown when Empty is true
    Verbose      bool       // Enable extra detail
    Wide         bool       // Enable wide format with additional columns
}
```

### MessageLevel

```go
type MessageLevel string

const (
    LevelInfo     MessageLevel = "info"
    LevelSuccess  MessageLevel = "success"
    LevelWarning  MessageLevel = "warning"
    LevelError    MessageLevel = "error"
    LevelDebug    MessageLevel = "debug"
    LevelProgress MessageLevel = "progress"
)
```

### Message

```go
type Message struct {
    Level   MessageLevel
    Content string
}
```

### Config

```go
type Config struct {
    Default      RendererName // Default renderer name
    Verbose      bool         // Enable verbose output
    NoColor      bool         // Disable colored output
    UseNerdFonts bool         // Enable Nerd Font icons
}

func DefaultConfig() Config
```

`DefaultConfig()` returns `Config{Default: RendererColored}`.

### Data Structures

Commands pass these types as the `data` argument to `Output`/`OutputWithContext`. Renderers handle each type appropriately.

#### TableData

```go
type TableData struct {
    Headers []string
    Rows    [][]string
}
```

#### ListData

```go
type ListData struct {
    Items []string
}
```

#### KeyValueData

```go
type KeyValueData struct {
    Pairs []KeyValue
}

type KeyValue struct {
    Key   string
    Value string
}
```

Constructors:

```go
// From a map (key order not guaranteed)
func NewKeyValueData(m map[string]string) KeyValueData

// With explicit ordering
func NewOrderedKeyValueData(pairs ...KeyValue) KeyValueData
```

---

## Renderer Interface

All six renderers implement this interface:

```go
type Renderer interface {
    Render(w io.Writer, data any, opts Options) error
    RenderWithContext(ctx context.Context, w io.Writer, data any, opts Options) error
    RenderMessage(w io.Writer, msg Message) error
    RenderMessageWithContext(ctx context.Context, w io.Writer, msg Message) error
    Name() RendererName
    SupportsColor() bool
}
```

`Render` and `RenderMessage` call their `WithContext` equivalents using `context.Background()` for backward compatibility.

---

## Renderers

### ColoredRenderer

The default renderer. Outputs richly formatted text with lipgloss styles and Unicode icons. When a `ColorProvider` is available in the context (via `colors.WithProvider`), it derives its styles from that provider; otherwise it falls back to built-in Catppuccin-inspired defaults.

```go
func NewColoredRenderer() *ColoredRenderer
func NewColoredRendererWithIcons(icons Icons) *ColoredRenderer
```

- `Name()` returns `RendererColored`.
- `SupportsColor()` returns `true`.
- Registered automatically via `init()`.

Handles: `KeyValueData`, `TableData`, `ListData`, `[]string`, `map[string]string`, `map[string]interface{}`, and any other type via `fmt.Fprintf`.

#### Icons

Three icon sets are provided:

```go
func DefaultIcons() Icons    // Unicode: ✓ ⚠ ✗ ℹ → • ▌
func NerdFontIcons() Icons   // Nerd Font codepoints
func PlainIcons() Icons      // ASCII: [OK] [!] [X] [i] -> * |
```

```go
type Icons struct {
    Success  string
    Warning  string
    Error    string
    Info     string
    Progress string
    Bullet   string
    Section  string
}
```

### PlainRenderer

Outputs plain text without color or styling. Suitable for piping, CI environments, and terminals without color support. Ignores context entirely.

```go
func NewPlainRenderer() *PlainRenderer
```

- `Name()` returns `RendererPlain`.
- `SupportsColor()` returns `false`.
- Registered automatically via `init()`.

Message prefixes: `[OK]`, `[WARN]`, `[ERROR]`, `[DEBUG]`, `->` (progress), `[INFO]`.

### JSONRenderer

Outputs data as indented JSON (2-space indent). Ignores `Title`, `EmptyMessage`, `EmptyHints`, and context. Converts `KeyValueData` to `map[string]string`, `TableData` to `[]map[string]string`, and `ListData` to `[]string`. Messages are output as `{"level": "...", "message": "..."}`.

```go
func NewJSONRenderer() *JSONRenderer
```

- `Name()` returns `RendererJSON`.
- `SupportsColor()` returns `false`.
- Registered automatically via `init()`.

### YAMLRenderer

Outputs data as YAML (2-space indent). Behavior mirrors `JSONRenderer` for data type conversions. Ignores `Title`, `EmptyMessage`, `EmptyHints`, and context.

```go
func NewYAMLRenderer() *YAMLRenderer
```

- `Name()` returns `RendererYAML`.
- `SupportsColor()` returns `false`.
- Registered automatically via `init()`.

### TableRenderer

Focuses exclusively on table output. Suppresses all `RenderMessage` calls (no-op). Skips output entirely when `opts.Empty` is `true`. Delegates actual table rendering to an embedded `ColoredRenderer`, so it does respond to context-provided colors.

```go
func NewTableRenderer() *TableRenderer
```

- `Name()` returns `RendererTable`.
- `SupportsColor()` returns `true`.
- Registered automatically via `init()`.

For `KeyValueData` and `map[string]string`, renders as `key: value` lines. For non-table data types, outputs nothing.

### CompactRenderer

Like `ColoredRenderer` but more condensed: tighter column spacing, muted-style headers (no separator line), `▸` title prefix instead of `▌`, and compact list items use `-` bullet instead of `•`. Embeds a `ColoredRenderer` and delegates to it for unsupported data types and for message rendering.

```go
func NewCompactRenderer() *CompactRenderer
```

- `Name()` returns `RendererCompact`.
- `SupportsColor()` returns `true`.
- Registered automatically via `init()`.

---

## Registry

The global registry maps `RendererName` to `Renderer` implementations. All six built-in renderers register themselves in their respective `init()` functions.

```go
func Register(r Renderer)
func Get(name RendererName) Renderer
func List() []RendererName
```

### Configuration

```go
func SetConfig(cfg Config)
func GetConfig() Config
func SetDefault(name RendererName)
```

### Writer

```go
func SetWriter(w io.Writer)
func GetWriter() io.Writer
```

The default writer is `os.Stdout`. Change it for testing by calling `SetWriter(buf)`.

### Renderer Resolution

```go
func ResolveRenderer(override string) Renderer
```

Resolves the renderer to use: override string > `DVM_RENDER` env var > global config default. If `NO_COLOR` is set and the resolved renderer is `colored` or `compact`, it falls back to `plain`. If the resolved renderer is not registered, falls back to `colored`, then `plain`.

---

## Output Functions

### Context-Aware (Preferred)

```go
func OutputWithContext(ctx context.Context, data any, opts Options) error
func OutputWithContextAndRenderer(ctx context.Context, rendererOverride string, data any, opts Options) error
func OutputToWithContext(ctx context.Context, w io.Writer, rendererOverride string, data any, opts Options) error
```

### Non-Context (Backward Compatible)

```go
func Output(data any, opts Options) error
func OutputWith(rendererOverride string, data any, opts Options) error
func OutputTo(w io.Writer, rendererOverride string, data any, opts Options) error
```

---

## Message Functions

### Context-Aware

```go
func MsgWithContext(ctx context.Context, level MessageLevel, content string) error
func MsgWithContextAndRenderer(ctx context.Context, rendererOverride string, level MessageLevel, content string) error
func MsgToWithContext(ctx context.Context, w io.Writer, rendererOverride string, msg Message) error
```

### Non-Context

```go
func Msg(level MessageLevel, content string) error
func MsgWith(rendererOverride string, level MessageLevel, content string) error
func MsgTo(w io.Writer, rendererOverride string, msg Message) error
```

### Convenience Message Functions

These call `Msg` with the corresponding `MessageLevel`:

```go
func Info(content string) error
func Success(content string) error
func Warning(content string) error
func Error(content string) error
func Progress(content string) error
```

---

## Convenience Functions

Formatted variants (use `fmt.Sprintf` internally):

```go
func Infof(format string, args ...any) error
func Successf(format string, args ...any) error
func Warningf(format string, args ...any) error
func Errorf(format string, args ...any) error
func Progressf(format string, args ...any) error
```

Stderr output:

```go
func InfoToStderr(content string) error
func WarningToStderr(content string) error
func ErrorToStderr(content string) error
func InfofToStderr(format string, args ...any) error
func WarningfToStderr(format string, args ...any) error
func ErrorfToStderr(format string, args ...any) error
```

Undecorated text output (no level prefix, no color):

```go
func Plain(text string) error
func Plainf(format string, args ...any) error
```

Empty line:

```go
func Blank() error
```
