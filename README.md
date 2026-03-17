# pf — Unified Port-Forward Manager

A CLI tool that manages port-forwarding to both Kubernetes services and EC2 instances from a single interface.

## Install

```bash
go install github.com/keimoon/pf@latest
```

Or build from source:

```bash
git clone https://github.com/keimoon/pf.git
cd pf && go install .
```

## Quick Start

```bash
# Set default SSH command and user (optional)
pf defaults --ssh-command ssh-nohost --user ec2-user

# Register services
pf add vm-us --type ec2 --host 198.51.100.1 -l 8428 -r 8428 --desc "VictoriaMetrics US"
pf add grafana --type k8s --target svc/grafana -n monitoring -l 3000 -r 80 --desc "Grafana"

# List registered services
pf list

# Forward a single service
pf connect vm-us

# Forward multiple services concurrently
pf connect vm-us grafana
# Ctrl+C stops all
```

## Commands

| Command | Description |
|---------|-------------|
| `pf add <name> [flags]` | Register a service |
| `pf remove <name>` | Unregister a service (alias: `rm`) |
| `pf list` | List all registered services (alias: `ls`) |
| `pf connect <name> [name...]` | Port-forward one or more services |
| `pf defaults [flags]` | View or set default configuration |
| `pf version` | Print version |

### `pf add` flags

| Flag | Description |
|------|-------------|
| `--type` | **Required.** `ec2` or `k8s` |
| `-l, --local` | **Required.** Local port |
| `-r, --remote` | **Required.** Remote port |
| `--host` | EC2: private IP or hostname |
| `--user` | EC2: SSH user (e.g. `ec2-user`, `ubuntu`) |
| `--target` | K8s: `svc/<name>`, `deploy/<name>`, or `pod/<name>` |
| `-n, --namespace` | K8s: namespace (default: `default`) |
| `--context` | K8s: kubectl context (default: current) |
| `--desc` | Description |
| `--ssh-command` | Per-service SSH command override |

### Global flags

| Flag | Description |
|------|-------------|
| `--config` | Config file path (default: `~/.config/pf/services.yaml`) |

## How It Works

- Services are stored in `~/.config/pf/services.yaml`, managed via `pf add` / `pf remove`
- `pf connect` delegates to `ssh -L` for EC2 or `kubectl port-forward` for K8s, running as child processes
- Multiple services are forwarded concurrently; Ctrl+C shuts them all down gracefully

## Config Example

```yaml
defaults:
    ssh_command: ssh-nohost
    ssh_user: ec2-user
services:
    grafana:
        type: k8s
        target: svc/grafana
        namespace: monitoring
        local_port: 3000
        remote_port: 80
        description: Grafana
    vm-us:
        type: ec2
        host: 198.51.100.1
        local_port: 8428
        remote_port: 8428
        description: VictoriaMetrics US
```
