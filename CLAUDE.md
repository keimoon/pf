# pf

Unified port-forward manager CLI for Kubernetes services and EC2 instances.

## Architecture

- **Entry point:** `main.go` calls `cmd.Execute()`
- **CLI framework:** Cobra (`cmd/` package) — one file per subcommand
- **Config:** `config/` package — YAML-backed at `~/.config/pf/services.yaml`, managed programmatically (not hand-edited)
- **Forwarding:** `forwarder/` package — `Forwarder` interface with `EC2Forwarder` (ssh -L) and `K8sForwarder` (kubectl port-forward) implementations

## Commands

- `cmd/root.go` — root command, global `--config` flag, `configPath()` helper
- `cmd/add.go` — `pf add <name>` registers a service
- `cmd/remove.go` — `pf remove <name>` unregisters a service
- `cmd/list.go` — `pf list` prints all services in a table
- `cmd/connect.go` — `pf connect <name>...` runs forwarders concurrently with signal handling
- `cmd/defaults.go` — `pf defaults` views/sets default SSH command

## Testing

```bash
go test ./...
```

Tests cover config validation, load/save round-trip, add/remove operations, SSH command fallback, and argument construction for both EC2 and K8s forwarders.

## Build

```bash
go build -o pf .
# or
go install .
```
