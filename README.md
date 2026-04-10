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

## Documentation

For comprehensive documentation, see the [MaestroSDK Documentation](https://rmkohlman.github.io/MaestroSDK/).

## Part of the DevOpsMaestro Ecosystem

MaestroSDK is a shared dependency of [DevOpsMaestro](https://github.com/rmkohlman/devopsmaestro), a kubectl-style CLI toolkit for containerized development environments.

## License

GPL-3.0 with commercial dual-license. See [LICENSE](LICENSE) for details.
