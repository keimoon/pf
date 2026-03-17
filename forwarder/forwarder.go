package forwarder

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/keimoon/pf/config"
)

type Forwarder interface {
	// Args returns the command and arguments for the port-forward.
	Args(svc config.Service) (string, []string)
}

// Run executes the forwarder, blocking until ctx is cancelled or the process exits.
func Run(ctx context.Context, f Forwarder, svc config.Service) error {
	bin, args := f.Args(svc)
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func New(cfg *config.Config, svc config.Service) (Forwarder, error) {
	switch svc.Type {
	case "ec2":
		return &EC2Forwarder{SSHCommand: cfg.SSHCommandFor(svc), SSHUser: cfg.SSHUserFor(svc)}, nil
	case "k8s":
		return &K8sForwarder{}, nil
	default:
		return nil, fmt.Errorf("unknown service type: %q", svc.Type)
	}
}
