# Yuppo

Automatic git commit watcher. Monitors a directory for file changes and creates commits with debouncing.

## Features

- **Auto-commit** — detects file changes (create, write, remove) and commits them automatically
- **Debouncing** — groups rapid changes into a single commit after a configurable quiet period
- **Gitignore support** — respects `.gitignore` patterns in the watched directory
- **Auto-push** — optionally pushes to origin after each commit
- **Recursive watching** — monitors all subdirectories automatically
- **Graceful shutdown** — handles SIGINT/SIGTERM cleanly

## Installation

```bash
go install github.com/sharki22/yuppo/cmd/yuppo@latest
```

Or build from source:

```bash
git clone https://github.com/sharki22/yuppo.git
cd yuppo
go build -o yuppo ./cmd/yuppo
```

## Configuration

Create a `config.yaml`:

```yaml
watch_path: /path/to/your/watch/folder/
commit_message: "auto: commit message"
auto_push: false
debounce_seconds: 20
```

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `watch_path` | string | yes | — | Directory to watch |
| `commit_message` | string | no | `"auto: update"` | Git commit message |
| `auto_push` | bool | no | `false` | Push to origin after commit |
| `debounce_seconds` | int | no | `5` | Seconds to wait before committing |

## Usage

```bash
# Default config path (config.yaml in current directory)
./yuppo

# Custom config path
./yuppo --config /path/to/config.yaml
```

## How it works

```
config.yaml
    │
    ▼
┌─────────┐    ┌──────────┐    ┌─────────┐
│  config  │───▶│  watcher │───▶│ gitops  │
└─────────┘    └──────────┘    └─────────┘
                    │               │
                    ▼               ▼
              ┌──────────┐   ┌──────────┐
              │gitignore │   │debouncer │
              └──────────┘   └──────────┘
```

1. **Config** loads and validates `config.yaml`
2. **Gitignore** loads `.gitignore` patterns from the watch path
3. **Watcher** recursively monitors the directory via `fsnotify`
4. **Debouncer** groups rapid file changes into a single event
5. **Gitops** runs `git add . && git commit` (and optionally `git push`)

## Project structure

```
yuppo/
├── cmd/yuppo/main.go          # entry point
├── internal/
│   ├── config/config.go       # config loading + validation
│   ├── debouncer/debouncer.go # event debouncing
│   ├── gitignore/gitignore.go # .gitignore parsing
│   ├── gitops/gitops.go       # git add/commit/push
│   └── watcher/watcher.go     # fsnotify watcher
├── config.yaml
├── go.mod
└── go.sum
```

## Future plans

- [ ] CLI flags (`--dry-run`, `--verbose`)
- [ ] Dry-run mode (preview commits without executing)
- [ ] Multiple watch paths
- [ ] Commit message templates with timestamps
- [ ] Webhook notifications (Discord, Telegram)
- [ ] Unit tests for all packages
- [ ] CI/CD with GitHub Actions
- [ ] Structured logging (`slog` / `zerolog`)
- [ ] systemd service file
- [ ] Config hot-reloading
- [ ] GPG commit signing support
