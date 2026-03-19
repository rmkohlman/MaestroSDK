# paths

Import path: `github.com/rmkohlman/MaestroSDK/paths`

The `paths` package provides a centralized, testable path configuration for all DevOpsMaestro tools (`dvm`, `nvp`, `dvt`). It replaces scattered hardcoded path constructions with a single `PathConfig` struct whose methods return deterministic paths derived from a home directory.

**Design goals:**

- No OS dependencies in tests — pass any home directory to `paths.New()`
- Single source of truth for all well-known filesystem locations
- Exported constants so callers can reference directory names without resolving full paths

---

## Constants

```go
const (
    DVMDirName   = ".devopsmaestro"
    NVPDirName   = ".nvp"
    DVTDirName   = ".dvt"
    DatabaseFile = "devopsmaestro.db"
)
```

| Constant | Value | Purpose |
|----------|-------|---------|
| `DVMDirName` | `".devopsmaestro"` | Hidden directory under `$HOME` for `dvm` state |
| `NVPDirName` | `".nvp"` | Hidden directory under `$HOME` for `nvp` state |
| `DVTDirName` | `".dvt"` | Hidden directory under `$HOME` for `dvt` state |
| `DatabaseFile` | `"devopsmaestro.db"` | SQLite database filename inside the `dvm` root |

---

## PathConfig

`PathConfig` is the central struct. It holds a resolved home directory and exposes methods that return fully-qualified paths for every well-known location in the DevOpsMaestro filesystem layout. The struct is immutable — `homeDir` is set once at construction and never changes.

### Constructors

#### `New`

```go
func New(homeDir string) *PathConfig
```

Creates a `PathConfig` rooted at the given home directory. Has no OS dependencies, making it ideal for tests.

Panics if `homeDir` is empty, because that indicates a programming error — every code path must supply a valid home directory.

```go
pc := paths.New("/tmp/fakehome")
```

#### `Default`

```go
func Default() (*PathConfig, error)
```

Creates a `PathConfig` using the current user's home directory returned by `os.UserHomeDir()`. This is the standard constructor for production code.

```go
pc, err := paths.Default()
if err != nil {
    // handle: unable to determine home directory
}
```

---

## Methods

### DVM Root

| Method | Returns | Example (home = `/home/user`) |
|--------|---------|-------------------------------|
| `Root()` | `{home}/.devopsmaestro` | `/home/user/.devopsmaestro` |
| `ConfigFile()` | `{root}/config.yaml` | `/home/user/.devopsmaestro/config.yaml` |
| `Database()` | `{root}/devopsmaestro.db` | `/home/user/.devopsmaestro/devopsmaestro.db` |
| `VersionFile()` | `{root}/.version` | `/home/user/.devopsmaestro/.version` |
| `ContextFile()` | `{root}/context.yaml` | `/home/user/.devopsmaestro/context.yaml` |
| `NvimSyncStatus()` | `{root}/.nvim-sync-status` | `/home/user/.devopsmaestro/.nvim-sync-status` |
| `LogsDir()` | `{root}/logs` | `/home/user/.devopsmaestro/logs` |
| `BackupsDir()` | `{root}/backups` | `/home/user/.devopsmaestro/backups` |
| `TemplatesDir()` | `{root}/templates` | `/home/user/.devopsmaestro/templates` |
| `NvimTemplatesDir()` | `{root}/templates/nvim` | `/home/user/.devopsmaestro/templates/nvim` |
| `ShellTemplatesDir()` | `{root}/templates/shell` | `/home/user/.devopsmaestro/templates/shell` |

### Workspace

| Method | Signature | Returns |
|--------|-----------|---------|
| `WorkspacesDir()` | `() string` | `{root}/workspaces` |
| `WorkspacePath()` | `(slug string) string` | `{root}/workspaces/{slug}` |
| `WorkspaceRepoPath()` | `(slug string) string` | `{root}/workspaces/{slug}/repo` |
| `WorkspaceVolumePath()` | `(slug string) string` | `{root}/workspaces/{slug}/volume` |
| `WorkspaceConfigPath()` | `(slug string) string` | `{root}/workspaces/{slug}/.dvm` |

### Git and Build

| Method | Signature | Returns |
|--------|-----------|---------|
| `ReposDir()` | `() string` | `{root}/repos` |
| `BuildStagingDir()` | `(appName string) string` | `{root}/build-staging/{appName}` |

### Registry

| Method | Signature | Returns |
|--------|-----------|---------|
| `RegistryDir()` | `(name string) string` | `{root}/registries/{name}` |
| `RegistryStorage()` | `() string` | `{root}/registry` |
| `AthensStorage()` | `() string` | `{root}/athens` |
| `VerdaccioStorage()` | `() string` | `{root}/verdaccio` |
| `DevpiStorage()` | `() string` | `{root}/devpi` |
| `SquidDir()` | `() string` | `{root}/squid` |

### NVP

| Method | Returns | Example (home = `/home/user`) |
|--------|---------|-------------------------------|
| `NVPRoot()` | `{home}/.nvp` | `/home/user/.nvp` |
| `NVPPluginsDir()` | `{nvpRoot}/plugins` | `/home/user/.nvp/plugins` |
| `NVPPackagesDir()` | `{nvpRoot}/packages` | `/home/user/.nvp/packages` |
| `NVPThemesDir()` | `{nvpRoot}/themes` | `/home/user/.nvp/themes` |
| `NVPCoreConfig()` | `{nvpRoot}/core.yaml` | `/home/user/.nvp/core.yaml` |

### DVT

| Method | Returns | Example (home = `/home/user`) |
|--------|---------|-------------------------------|
| `DVTRoot()` | `{home}/.dvt` | `/home/user/.dvt` |
| `DVTPromptsDir()` | `{dvtRoot}/prompts` | `/home/user/.dvt/prompts` |
| `DVTPluginsDir()` | `{dvtRoot}/plugins` | `/home/user/.dvt/plugins` |
| `DVTShellsDir()` | `{dvtRoot}/shells` | `/home/user/.dvt/shells` |
| `DVTProfilesDir()` | `{dvtRoot}/profiles` | `/home/user/.dvt/profiles` |
| `DVTActiveProfile()` | `{dvtRoot}/.active-profile` | `/home/user/.dvt/.active-profile` |

### Helper

#### `DatabasePathTilde`

```go
func (p *PathConfig) DatabasePathTilde() string
```

Returns the tilde-notation string `~/.devopsmaestro/devopsmaestro.db`. This is **not** a real filesystem path — it is intended as a default config value in viper configurations for `nvp` and `dvt`, which expand the tilde at runtime.

```go
pc, _ := paths.Default()
viperDefault := pc.DatabasePathTilde() // "~/.devopsmaestro/devopsmaestro.db"
```

---

## Usage Examples

### Production code

```go
import "github.com/rmkohlman/MaestroSDK/paths"

pc, err := paths.Default()
if err != nil {
    return fmt.Errorf("cannot resolve paths: %w", err)
}

dbPath  := pc.Database()             // e.g. /home/user/.devopsmaestro/devopsmaestro.db
wsPath  := pc.WorkspacePath("dev")  // e.g. /home/user/.devopsmaestro/workspaces/dev
nvpRoot := pc.NVPRoot()              // e.g. /home/user/.nvp
```

### Tests

```go
import "github.com/rmkohlman/MaestroSDK/paths"

pc := paths.New("/tmp/fakehome")

dbPath := pc.Database()
// "/tmp/fakehome/.devopsmaestro/devopsmaestro.db" — no OS dependency
```

### Using constants without a full path

```go
import "github.com/rmkohlman/MaestroSDK/paths"

// Skip the devopsmaestro directory when walking the filesystem
if entry.Name() == paths.DVMDirName {
    return filepath.SkipDir
}
```
