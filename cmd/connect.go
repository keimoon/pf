package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/keimoon/pf/config"
	"github.com/keimoon/pf/forwarder"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect <name> [name...]",
	Short: "Port-forward one or more services",
	Long: `Start port-forwarding for the named services. Examples:

  pf connect vm-us
  pf connect vm-us grafana    # forward multiple services concurrently`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configPath()
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		// Validate all names before starting any forwards
		for _, name := range args {
			if _, ok := cfg.Services[name]; !ok {
				return fmt.Errorf("service %q not found (run 'pf list' to see available services)", name)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle Ctrl+C gracefully
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			fmt.Fprintln(os.Stderr, "\nShutting down...")
			cancel()
		}()

		var wg sync.WaitGroup
		errCh := make(chan error, len(args))

		for _, name := range args {
			svc := cfg.Services[name]
			fwd, err := forwarder.New(cfg, svc)
			if err != nil {
				cancel()
				return err
			}

			target := svc.Host
			if svc.Type == "k8s" {
				target = svc.Target
			}
			fmt.Printf("%-12s  %s localhost:%d -> %s:%d\n",
				name, svc.Type, svc.LocalPort, target, svc.RemotePort)

			wg.Add(1)
			go func(name string, fwd forwarder.Forwarder, svc config.Service) {
				defer wg.Done()
				if err := forwarder.Run(ctx, fwd, svc); err != nil && ctx.Err() == nil {
					errCh <- fmt.Errorf("%s: %w", name, err)
				}
			}(name, fwd, svc)
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
