# MaestroSDK

Shared SDK packages for the DevOpsMaestro ecosystem.

## What is MaestroSDK?

MaestroSDK is a standalone Go module (`github.com/rmkohlman/MaestroSDK`) that provides shared infrastructure packages used across the DevOpsMaestro toolchain. It is designed to be imported by `dvm`, `nvp`, `dvt`, and any other tools in the ecosystem that need common color, rendering, resource, or path abstractions.

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

Inject a theme-aware `ColorProvider` into a `context.Context` and retrieve it anywhere downstream. Call `colors.InitColorProviderForCommand` in a CLI command's setup to wire up the provider (respecting `--no-color` and the `NO_COLOR` env var), then use `colors.FromContextOrDefault` anywhere downstream to access colors.

### render

Prepare structured data in your command (e.g., `render.TableData` with headers and rows) and pass it to `render.OutputWithContext`. The renderer is selected by the `DVM_RENDER` environment variable or an explicit override flag.

### resource

Register a handler once at startup, then use the package-level `Apply`, `Get`, `List`, and `Delete` functions to manage resources by kind and name.

### paths

Call `paths.Default()` to get a `PathConfig` rooted at the user's home directory. All well-known DevOpsMaestro filesystem locations are available as methods: `Database()`, `WorkspacePath(slug)`, `NVPRoot()`, and more. In tests, use `paths.New("/tmp/fakehome")` for a fully deterministic, OS-independent path config.

## Part of the DevOpsMaestro Ecosystem

MaestroSDK is a shared dependency of [DevOpsMaestro](https://github.com/rmkohlman/devopsmaestro), a kubectl-style CLI toolkit for containerized development environments.

## License

GPL-3.0 with commercial dual-license. See [LICENSE](https://github.com/rmkohlman/MaestroSDK/blob/main/LICENSE) for details.
