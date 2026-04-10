# resource

The `resource` package provides a unified, kubectl-style interface for managing resources in DevOpsMaestro. It follows the pattern where resources are identified by `Kind` and are applied, retrieved, listed, and deleted through a common interface.

**Architecture:**

- `Resource`: The data being managed (plugins, themes, workspaces, etc.)
- `Handler`: Knows how to CRUD a specific resource type
- Registry: Routes operations to the correct handler by `Kind`

---

## Core Types

### Resource Interface

Any type that exposes `GetKind()`, `GetName()`, and `Validate()` satisfies the `Resource` interface.

| Method | Returns | Description |
|--------|---------|-------------|
| `GetKind()` | `string` | Resource type (e.g., `"NvimPlugin"`, `"Workspace"`) |
| `GetName()` | `string` | Unique name of this resource |
| `Validate()` | `error` | Validates the resource's fields |

### Handler Interface

Each resource type has a corresponding handler:

| Method | Description |
|--------|-------------|
| `Kind()` | Resource type this handler manages |
| `Apply(ctx, data)` | Create or update from YAML bytes |
| `Get(ctx, name)` | Retrieve by name |
| `List(ctx)` | List all of this type |
| `Delete(ctx, name)` | Remove by name |
| `ToYAML(res)` | Serialize to YAML bytes |

### Context

`Context` provides dependencies needed by handlers without tight coupling to any specific implementation:

| Field | Type | Description |
|-------|------|-------------|
| `DataStore` | `any` | Database store; type-assert to the specific store interface |
| `PluginStore` | `any` | Pre-configured plugin store (optional) |
| `ThemeStore` | `any` | Pre-configured theme store (optional) |
| `ConfigDir` | `string` | Configuration directory for file-based storage |

`DataStore`, `PluginStore`, and `ThemeStore` are typed as `any`. Handlers use the `DataStoreAs`, `PluginStoreAs`, and `ThemeStoreAs` generic helpers to safely assert these to their required interface types.

### KindHeader

Used internally to detect the `Kind` field from YAML before full parsing.

### DetectKind

Extracts the `kind` field from YAML without fully parsing the document. Returns an error if YAML is malformed or the `kind` field is absent.

---

## Registry

The registry maps `Kind` strings to `Handler` implementations. It is safe for concurrent use.

### Registration

| Function | Description |
|----------|-------------|
| `Register(h Handler)` | Adds a handler. Panics if a handler for the same `Kind` is already registered. Use in startup initialization. |
| `RegisterSafe(h Handler)` | Like `Register` but returns an error on duplicate instead of panicking. Prefer this outside startup. |
| `SetFallbackHandler(h Handler)` | Sets a handler for unknown `Kind` values (e.g., a dynamic/custom resource handler). |

### Lookup

| Function | Description |
|----------|-------------|
| `GetHandler(kind)` | Returns the handler for `kind`, or the fallback handler if none is registered, or `nil` if no fallback is set. |
| `MustGetHandler(kind)` | Returns the handler or an error if none is registered. |
| `RegisteredKinds()` | Returns all registered `Kind` strings. |
| `ClearRegistry()` | Removes all registered handlers. Intended for testing only. |

### Package-Level CRUD Operations

These functions look up the appropriate handler and delegate to it:

| Function | Description |
|----------|-------------|
| `Apply(ctx, data, source)` | Parses YAML, calls `DetectKind` to find the handler, applies the resource |
| `Get(ctx, kind, name)` | Retrieves one resource by kind and name |
| `List(ctx, kind)` | Lists all resources of a kind |
| `Delete(ctx, kind, name)` | Deletes a resource by kind and name |
| `ToYAML(res)` | Serializes a resource to YAML |

---

## ResourceList

`ResourceList` is a kubectl-style `List` wrapper for exporting and importing multiple resources in a single YAML document. It is not stored in the database; it is produced by export commands and consumed by `apply`.

A `ResourceList` has `apiVersion: devopsmaestro.io/v1` and `kind: List`.

| Function | Description |
|----------|-------------|
| `NewResourceList()` | Creates an empty `ResourceList` |
| `BuildList(ctx, resources)` | Serializes each `Resource` via its handler's `ToYAML()` and adds it to the list's `Items`. Resources that fail serialization are skipped with a warning. |
| `ApplyList(ctx, data)` | Parses a `List` YAML document and applies each item. Continues on error (kubectl precedent). Returns all successfully applied resources and a summary error if any items failed. |

### DependencyOrder

`DependencyOrder` defines the canonical order in which resource kinds should be applied so that dependencies come before dependents:

`Ecosystem` → `Domain` → `App` → `GitRepo` → `Registry` → `Credential` → `Workspace` → `NvimPlugin` → `NvimTheme` → `NvimPackage` → `TerminalPrompt` → `TerminalPackage`

---

## DataStore Helpers

These generic functions provide type-safe extraction of typed store interfaces from a `Context`, without coupling the `resource` package to any specific store implementation.

| Function | Description |
|----------|-------------|
| `DataStoreAs[T](ctx)` | Extracts `ctx.DataStore` and asserts it to type `T`. Returns an error if `DataStore` is `nil` or if the type assertion fails. |
| `PluginStoreAs[T](ctx)` | Extracts `ctx.PluginStore` and asserts it to type `T`. |
| `ThemeStoreAs[T](ctx)` | Extracts `ctx.ThemeStore` and asserts it to type `T`. |
