# MaestroSDK

Shared SDK packages for the DevOpsMaestro ecosystem.

## What is MaestroSDK?

MaestroSDK is a standalone Go module (`github.com/rmkohlman/MaestroSDK`) that provides shared infrastructure packages used across the DevOpsMaestro toolchain. It is designed to be imported by `dvm`, `nvp`, `dvt`, and any other tools in the ecosystem that need common color, rendering, resource, or path abstractions.

## Installation

```bash
go get github.com/rmkohlman/MaestroSDK
```

**Go version requirement:** 1.25.6 or later.

## Packages

| Package | Import Path | Description |
|---------|-------------|-------------|
| `colors` | `github.com/rmkohlman/MaestroSDK/colors` | `ColorProvider` interface, theme-aware color injection via `context.Context`, factory, CLI helpers, and mock for testing |
| `render` | `github.com/rmkohlman/MaestroSDK/render` | Output rendering system with six renderers: colored, plain, JSON, YAML, table, and compact |
| `resource` | `github.com/rmkohlman/MaestroSDK/resource` | kubectl-style resource handler registry, `ResourceList` with dependency ordering, generic DataStore helpers |
| `paths` | `github.com/rmkohlman/MaestroSDK/paths` | Centralized, testable path configuration for all well-known DevOpsMaestro filesystem locations |

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/charmbracelet/lipgloss` | Terminal styling in the `render` package |
| `github.com/rmkohlman/MaestroPalette` | Palette data model consumed by the `colors` package |
| `gopkg.in/yaml.v3` | YAML encoding in `render` and `resource` packages |
| `github.com/stretchr/testify` | Test assertions |

## Quick Start

### colors

Inject a theme-aware `ColorProvider` into a `context.Context` and retrieve it anywhere downstream:

```go
import "github.com/rmkohlman/MaestroSDK/colors"

// In a CLI command (e.g. cobra PersistentPreRunE)
ctx, err := colors.InitColorProviderForCommand(ctx, paletteProvider, noColorFlag)

// Anywhere downstream
provider := colors.FromContextOrDefault(ctx)
successColor := provider.Success()  // e.g. "#9ece6a"
```

### render

Prepare structured data in your command and pass it to the global `Output` function:

```go
import "github.com/rmkohlman/MaestroSDK/render"

data := render.TableData{
    Headers: []string{"NAME", "STATUS", "APP"},
    Rows: [][]string{
        {"dev", "running", "myapp"},
        {"staging", "stopped", "myapp"},
    },
}

err := render.OutputWithContext(ctx, data, render.Options{
    Type:  render.TypeTable,
    Title: "Workspaces",
})
```

The renderer is selected by: `DVM_RENDER` environment variable, or the `-r`/`--render` flag value passed to `OutputWithContextAndRenderer`.

### resource

Register a handler once at startup, then use package-level functions to apply, get, list, and delete resources:

```go
import "github.com/rmkohlman/MaestroSDK/resource"

// Registration (typically in an init() or startup function)
resource.Register(&MyHandler{})

// Apply YAML from any source
res, err := resource.Apply(ctx, yamlData, "source.yaml")

// List all of a kind
items, err := resource.List(ctx, "MyKind")
```

### paths

```go
import "github.com/rmkohlman/MaestroSDK/paths"

pc, err := paths.Default()           // uses os.UserHomeDir()
// pc := paths.New("/tmp/fakehome")  // deterministic alternative for tests

dbPath  := pc.Database()             // ~/.devopsmaestro/devopsmaestro.db
wsPath  := pc.WorkspacePath("myws") // ~/.devopsmaestro/workspaces/myws
nvpRoot := pc.NVPRoot()              // ~/.nvp
```

## Part of the DevOpsMaestro Ecosystem

MaestroSDK is a shared dependency of [DevOpsMaestro](https://github.com/rmkohlman/devopsmaestro), a kubectl-style CLI toolkit for containerized development environments.

## License

GPL-3.0 with commercial dual-license. See [LICENSE](https://github.com/rmkohlman/MaestroSDK/blob/main/LICENSE) for details.
