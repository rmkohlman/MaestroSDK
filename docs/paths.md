# paths

The `paths` package provides a centralized, testable path configuration for all DevOpsMaestro tools (`dvm`, `nvp`, `dvt`). It replaces scattered hardcoded path constructions with a single `PathConfig` whose methods return deterministic paths derived from a home directory.

**Design goals:**

- No OS dependencies in tests — pass any home directory to `paths.New()`
- Single source of truth for all well-known filesystem locations
- Exported constants so callers can reference directory names without resolving full paths

---

## Constants

| Constant | Value | Purpose |
|----------|-------|---------|
| `DVMDirName` | `".devopsmaestro"` | Hidden directory under `$HOME` for `dvm` state |
| `NVPDirName` | `".nvp"` | Hidden directory under `$HOME` for `nvp` state |
| `DVTDirName` | `".dvt"` | Hidden directory under `$HOME` for `dvt` state |
| `DatabaseFile` | `"devopsmaestro.db"` | SQLite database filename inside the `dvm` root |

---

## PathConfig

`PathConfig` is the central type. It holds a resolved home directory and exposes methods that return fully-qualified paths for every well-known location in the DevOpsMaestro filesystem layout. The struct is immutable — `homeDir` is set once at construction and never changes.

### Constructors

**`New(homeDir string)`** — Creates a `PathConfig` rooted at the given home directory. Has no OS dependencies, making it ideal for tests. Panics if `homeDir` is empty, because that indicates a programming error.

**`Default()`** — Creates a `PathConfig` using the current user's home directory. This is the standard constructor for production code.

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

| Method | Returns |
|--------|---------|
| `WorkspacesDir()` | `{root}/workspaces` |
| `WorkspacePath(slug)` | `{root}/workspaces/{slug}` |
| `WorkspaceRepoPath(slug)` | `{root}/workspaces/{slug}/repo` |
| `WorkspaceVolumePath(slug)` | `{root}/workspaces/{slug}/volume` |
| `WorkspaceConfigPath(slug)` | `{root}/workspaces/{slug}/.dvm` |

### Git and Build

| Method | Returns |
|--------|---------|
| `ReposDir()` | `{root}/repos` |
| `BuildStagingDir(appName)` | `{root}/build-staging/{appName}` |

### Registry

| Method | Returns |
|--------|---------|
| `RegistryDir(name)` | `{root}/registries/{name}` |
| `RegistryStorage()` | `{root}/registry` |
| `AthensStorage()` | `{root}/athens` |
| `VerdaccioStorage()` | `{root}/verdaccio` |
| `DevpiStorage()` | `{root}/devpi` |
| `SquidDir()` | `{root}/squid` |

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

**`DatabasePathTilde()`** — Returns the tilde-notation string `~/.devopsmaestro/devopsmaestro.db`. This is **not** a real filesystem path — it is intended as a default config value for tools that expand the tilde at runtime.
