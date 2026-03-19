# resource

Import path: `github.com/rmkohlman/MaestroSDK/resource`

The `resource` package provides a unified, kubectl-style interface for managing resources in DevOpsMaestro. It follows the pattern where resources are identified by `Kind` and are applied, retrieved, listed, and deleted through a common interface.

**Architecture:**

- `Resource`: The data being managed (plugins, themes, workspaces, etc.)
- `Handler`: Knows how to CRUD a specific resource type
- Registry: Routes operations to the correct handler by `Kind`

---

## Core Types

### Resource Interface

```go
type Resource interface {
    GetKind() string  // Resource type (e.g., "NvimPlugin", "Workspace")
    GetName() string  // Unique name of this resource
    Validate() error  // Validates the resource's fields
}
```

### Handler Interface

Each resource type has a corresponding handler that implements this interface:

```go
type Handler interface {
    Kind() string                                    // Resource type this handler manages
    Apply(ctx Context, data []byte) (Resource, error) // Create or update from YAML
    Get(ctx Context, name string) (Resource, error)   // Retrieve by name
    List(ctx Context) ([]Resource, error)             // List all of this type
    Delete(ctx Context, name string) error            // Remove by name
    ToYAML(res Resource) ([]byte, error)              // Serialize to YAML bytes
}
```

### Context

`Context` provides dependencies needed by handlers without tight coupling to any specific implementation:

```go
type Context struct {
    DataStore   any    // Database store; type-assert to the specific store interface
    PluginStore any    // Pre-configured plugin store (optional)
    ThemeStore  any    // Pre-configured theme store (optional)
    ConfigDir   string // Configuration directory for file-based storage (e.g., ~/.nvp)
}
```

`DataStore`, `PluginStore`, and `ThemeStore` are typed as `any`. Handlers use the `DataStoreAs`, `PluginStoreAs`, and `ThemeStoreAs` generic helpers to safely assert these to their required interface types.

### KindHeader

Used internally to detect the `Kind` field from YAML before full parsing:

```go
type KindHeader struct {
    Kind string `yaml:"kind"`
}
```

### DetectKind

```go
func DetectKind(data []byte) (string, error)
```

Extracts the `kind` field from YAML without fully parsing the document. Returns an error if YAML is malformed or the `kind` field is absent.

---

## Registry

The registry maps `Kind` strings to `Handler` implementations. It is safe for concurrent use.

### Registration

```go
func Register(h Handler)
```

Adds a handler. Panics if a handler for the same `Kind` is already registered. Use in `init()` functions.

```go
func RegisterSafe(h Handler) error
```

Like `Register` but returns an error on duplicate instead of panicking. Prefer this outside `init()`.

```go
func SetFallbackHandler(h Handler)
```

Sets a handler for unknown `Kind` values (e.g., a dynamic/custom resource handler).

### Lookup

```go
func GetHandler(kind string) Handler
```

Returns the handler for `kind`, or the fallback handler if none is registered for that kind, or `nil` if no fallback is set.

```go
func MustGetHandler(kind string) (Handler, error)
```

Returns the handler or an error (`"no handler registered for kind: <kind>"`).

```go
func RegisteredKinds() []string
```

Returns all registered `Kind` strings.

```go
func ClearRegistry()
```

Removes all registered handlers and clears the fallback. Intended for testing only.

### Package-Level CRUD Operations

These functions look up the appropriate handler and delegate to it:

```go
func Apply(ctx Context, data []byte, source string) (Resource, error)
func Get(ctx Context, kind, name string) (Resource, error)
func List(ctx Context, kind string) ([]Resource, error)
func Delete(ctx Context, kind, name string) error
func ToYAML(res Resource) ([]byte, error)
```

`Apply` calls `DetectKind` on the data to determine which handler to invoke.

---

## ResourceList

`ResourceList` is a kubectl-style `List` wrapper for exporting and importing multiple resources in a single YAML document. It is not stored in the database; it is produced by export commands and consumed by `apply`.

```go
type ResourceList struct {
    APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
    Kind       string         `json:"kind"       yaml:"kind"`
    Metadata   map[string]any `json:"metadata"   yaml:"metadata"`
    Items      []any          `json:"items"      yaml:"items"`
}
```

### Constructors and Operations

```go
func NewResourceList() *ResourceList
```

Creates an empty `ResourceList` with `APIVersion: "devopsmaestro.io/v1"` and `Kind: "List"`.

```go
func BuildList(ctx Context, resources []Resource) (*ResourceList, error)
```

Serializes each `Resource` via its handler's `ToYAML()` method and adds it to the list's `Items`. Resources that fail serialization are skipped with a `slog.Warn` call; they do not cause the function to return an error. The caller is responsible for passing resources in the desired order.

```go
func ApplyList(ctx Context, data []byte) ([]Resource, error)
```

Parses a `List` YAML document and applies each item via the registered handler for its `Kind`. Continues on error (kubectl precedent). Returns all successfully applied resources and a summary error if any items failed (`"N of M items failed to apply"`).

### DependencyOrder

`DependencyOrder` defines the canonical order in which resource kinds should be applied so that dependencies come before dependents:

```go
var DependencyOrder = []string{
    "Ecosystem",
    "Domain",
    "App",
    "GitRepo",
    "Registry",
    "Credential",
    "Workspace",
    "NvimPlugin",
    "NvimTheme",
    "NvimPackage",
    "TerminalPrompt",
    "TerminalPackage",
}
```

---

## DataStore Helpers

These generic functions provide compile-time-safe extraction of typed store interfaces from a `Context`, without coupling the `resource` package to any specific store implementation.

### DataStoreAs

```go
func DataStoreAs[T any](ctx Context) (T, error)
```

Extracts `ctx.DataStore` and asserts it to type `T`. Returns an error if `DataStore` is `nil` or if the type assertion fails.

```go
// Example
ds, err := resource.DataStoreAs[db.DataStore](ctx)
```

### PluginStoreAs

```go
func PluginStoreAs[T any](ctx Context) (T, error)
```

Extracts `ctx.PluginStore` and asserts it to type `T`.

```go
// Example
ps, err := resource.PluginStoreAs[store.PluginStore](ctx)
```

### ThemeStoreAs

```go
func ThemeStoreAs[T any](ctx Context) (T, error)
```

Extracts `ctx.ThemeStore` and asserts it to type `T`.

```go
// Example
ts, err := resource.ThemeStoreAs[theme.Store](ctx)
```

---

## Usage Example

```go
// 1. Implement Handler for your resource type
type MyHandler struct{}

func (h *MyHandler) Kind() string { return "MyKind" }

func (h *MyHandler) Apply(ctx resource.Context, data []byte) (resource.Resource, error) {
    ds, err := resource.DataStoreAs[MyDataStore](ctx)
    if err != nil {
        return nil, err
    }
    // ... parse data, call ds.Upsert(...)
}
// ... implement Get, List, Delete, ToYAML

// 2. Register at startup
func init() {
    resource.Register(&MyHandler{})
}

// 3. Use package-level functions
ctx := resource.Context{DataStore: myDS}
res, err := resource.Apply(ctx, yamlBytes, "config.yaml")
items, err := resource.List(ctx, "MyKind")
err = resource.Delete(ctx, "MyKind", "item-name")
```
