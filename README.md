# MaestroSDK

Shared SDK packages for the DevOpsMaestro ecosystem.

## Overview

MaestroSDK (`github.com/rmkohlman/MaestroSDK`) provides four standalone packages used across the DevOpsMaestro toolchain:

| Package | Description |
|---------|-------------|
| `colors` | `ColorProvider` interface, theme-aware color injection via context, factory and CLI helpers |
| `render` | Output rendering system with six renderers (colored, plain, JSON, YAML, table, compact) |
| `resource` | kubectl-style resource handler registry, `ResourceList` with dependency ordering, generic DataStore helpers |
| `paths` | Centralized, testable path configuration for all DevOpsMaestro filesystem locations |

## Installation

```bash
go get github.com/rmkohlman/MaestroSDK
```

## Quick Examples

### colors -- inject and retrieve theme colors via context

```go
import "github.com/rmkohlman/MaestroSDK/colors"

// Initialize for a CLI command (respects --no-color and NO_COLOR env var)
ctx, err := colors.InitColorProviderForCommand(ctx, paletteProvider, noColorFlag)

// Retrieve the provider anywhere downstream
provider := colors.FromContextOrDefault(ctx)
fmt.Println(provider.Primary())   // e.g. "#7aa2f7"
fmt.Println(provider.Success())   // e.g. "#9ece6a"
```

### render -- output structured data in any format

```go
import "github.com/rmkohlman/MaestroSDK/render"

data := render.TableData{
    Headers: []string{"NAME", "STATUS"},
    Rows:    [][]string{{"myapp", "running"}},
}

// Uses the default colored renderer; override with DVM_RENDER env var or -r flag
err := render.Output(data, render.Options{
    Type:  render.TypeTable,
    Title: "Workspaces",
})
```

### paths -- deterministic filesystem paths

```go
import "github.com/rmkohlman/MaestroSDK/paths"

// Production: resolves from os.UserHomeDir()
pc, err := paths.Default()

// Tests: fully deterministic, no OS dependency
pc := paths.New("/tmp/fakehome")

dbPath := pc.Database()            // /tmp/fakehome/.devopsmaestro/devopsmaestro.db
ws     := pc.WorkspacePath("myws") // /tmp/fakehome/.devopsmaestro/workspaces/myws
```

## Documentation

For comprehensive documentation, see the [MaestroSDK Documentation](https://rmkohlman.github.io/MaestroSDK/).

## Part of the DevOpsMaestro Ecosystem

MaestroSDK is a shared dependency of [DevOpsMaestro](https://github.com/rmkohlman/devopsmaestro), a kubectl-style CLI toolkit for containerized development environments.

## License

GPL-3.0 with commercial dual-license. See [LICENSE](LICENSE) for details.
