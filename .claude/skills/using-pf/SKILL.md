---
name: using-pf
description: Use when the user wants to port-forward to EC2 instances or Kubernetes services using the `pf` CLI tool. Trigger this skill whenever the user mentions port-forwarding, SSH tunnels to EC2 hosts, kubectl port-forward, connecting to remote services, or managing their forwarding config — even if they don't say "pf" by name. Also use when they reference `~/.config/pf/services.yaml` or ask about registering/connecting/listing services for forwarding.
---

# Using `pf` — Unified Port-Forward Manager

`pf` manages port-forwarding to Kubernetes services and EC2 instances from a single CLI. Services are registered once, then connected by name.

## Core Workflow

```bash
# 1. (One-time) Set default SSH command and user if needed
pf defaults --ssh-command ssh-nohost --user ec2-user

# 2. Register services
pf add grafana --type k8s --target svc/grafana -n monitoring -l 3000 -r 80 --desc "Grafana"
pf add vm-us --type ec2 --host 198.51.100.1 -l 8428 -r 8428 --desc "VictoriaMetrics US"

# 3. Connect (Ctrl+C to stop)
pf connect grafana vm-us
```

## Commands

### `pf add <name> --type ec2|k8s [flags]`

Register a service. The name is how you'll refer to it in `connect`, `list`, and `remove`.

**Required flags** (all service types):
| Flag | Description |
|------|-------------|
| `--type` | `ec2` or `k8s` |
| `-l, --local` | Local port to listen on |
| `-r, --remote` | Remote port to forward to |

**EC2-specific:**
| Flag | Description |
|------|-------------|
| `--host` | **Required for EC2.** Private IP or hostname |
| `--user` | SSH user (e.g. `ec2-user`, `ubuntu`) — overrides global default |
| `--ssh-command` | Override SSH binary for this service only |

**K8s-specific:**
| Flag | Description |
|------|-------------|
| `--target` | **Required for K8s.** Resource to forward: `svc/<name>`, `deploy/<name>`, or `pod/<name>` |
| `-n, --namespace` | Namespace (default: `default`) |
| `--context` | kubectl context (default: current) |

**Optional:**
| Flag | Description |
|------|-------------|
| `--desc` | Human-readable description shown in `pf list` |

**Examples:**

```bash
# EC2 instance with custom SSH and user
pf add vm-eu --type ec2 --host 198.51.100.2 --user ubuntu -l 8429 -r 8428 --desc "VictoriaMetrics EU" --ssh-command ssh

# K8s service in a specific context
pf add staging-db --type k8s --target svc/postgres -n database --context staging -l 5432 -r 5432 --desc "Staging Postgres"

# K8s pod directly
pf add debug-pod --type k8s --target pod/my-debug-pod -l 9090 -r 9090
```

Duplicate names are rejected — remove the old one first with `pf remove`.

### `pf remove <name>` (alias: `pf rm`)

Unregister a service by name.

```bash
pf remove vm-eu
pf rm debug-pod     # alias works too
```

### `pf list` (alias: `pf ls`)

Show all registered services in a table:

```
NAME        TYPE  TARGET                  LOCAL  REMOTE  DESCRIPTION
grafana     k8s   monitoring/svc/grafana  3000   80      Grafana
vm-us       ec2   198.51.100.1           8428   8428    VictoriaMetrics US
```

### `pf connect <name> [name...]`

Start port-forwarding. Pass one or multiple names — they run concurrently. Ctrl+C stops all of them.

```bash
pf connect vm-us                 # single service
pf connect vm-us grafana vm-eu   # all three at once
```

All names are validated before any forwarding starts, so a typo won't leave you with a partial connection.

Under the hood:
- **EC2** services run `<ssh-command> -N -L <local>:localhost:<remote> [user@]<host>`
- **K8s** services run `kubectl port-forward [-n <ns>] [--context <ctx>] <target> <local>:<remote>`

### `pf defaults [--ssh-command <cmd>] [--user <user>]`

View or set global defaults. Without flags, shows current values. With flags, updates them.

```bash
pf defaults                              # show current defaults
pf defaults --ssh-command ssh-nohost     # set SSH command
pf defaults --user ec2-user              # set SSH user
pf defaults --ssh-command ssh --user root # set both at once
```

Resolution order for both SSH command and user: per-service override > global default > system default (`ssh` / current OS user).

### `pf version`

Print the version of `pf`.

### Global flag: `--config`

All commands accept `--config <path>` to use an alternative config file instead of `~/.config/pf/services.yaml`.

## Config File

Stored at `~/.config/pf/services.yaml`. Managed entirely through `pf add`, `pf remove`, and `pf defaults` — no need to hand-edit. Structure:

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

## Common Scenarios

**"I need to access a service running on an EC2 instance"**
→ `pf add <name> --type ec2 --host <ip> -l <port> -r <port>` then `pf connect <name>`

**"I need to access a K8s service locally"**
→ `pf add <name> --type k8s --target svc/<svc-name> -n <namespace> -l <local-port> -r <svc-port>` then `pf connect <name>`

**"I want to forward several things at once for my dev environment"**
→ Register each service with `pf add`, then `pf connect svc1 svc2 svc3`

**"Port is already in use"**
→ Pick a different local port with `-l`. The remote port stays the same.

**"Service name already exists"**
→ `pf remove <name>` first, then `pf add` again.

## Development

The project is a Go CLI built with Cobra:

```
main.go           → entry point, calls cmd.Execute()
cmd/              → one file per subcommand (root, add, remove, list, connect, defaults, version)
config/           → Config types, Load/Save/Add/Remove, validation
forwarder/        → Forwarder interface, EC2Forwarder (ssh), K8sForwarder (kubectl)
```

Run tests: `go test ./...`
Build: `go build -o pf .`
Install: `go install .`
