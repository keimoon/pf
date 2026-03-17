package cmd

import (
	"fmt"

	"github.com/keimoon/pf/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Register a new service for port-forwarding",
	Long: `Register a new service. Examples:

  pf add vm-us --type ec2 --host 198.51.100.1 -l 8428 -r 8428 --desc "VictoriaMetrics US"
  pf add grafana --type k8s --target svc/grafana -n monitoring -l 3000 -r 80 --desc "Grafana"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		svcType, _ := cmd.Flags().GetString("type")
		host, _ := cmd.Flags().GetString("host")
		target, _ := cmd.Flags().GetString("target")
		namespace, _ := cmd.Flags().GetString("namespace")
		context, _ := cmd.Flags().GetString("context")
		localPort, _ := cmd.Flags().GetInt("local")
		remotePort, _ := cmd.Flags().GetInt("remote")
		desc, _ := cmd.Flags().GetString("desc")
		sshCmd, _ := cmd.Flags().GetString("ssh-command")
		sshUser, _ := cmd.Flags().GetString("user")

		svc := config.Service{
			Type:        svcType,
			Host:        host,
			Target:      target,
			Namespace:   namespace,
			Context:     context,
			LocalPort:   localPort,
			RemotePort:  remotePort,
			Description: desc,
			SSHCommand:  sshCmd,
			SSHUser:     sshUser,
		}

		path := configPath()
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		if err := cfg.Add(name, svc); err != nil {
			return err
		}

		if err := config.Save(path, cfg); err != nil {
			return err
		}

		fmt.Printf("Added %q (%s, localhost:%d -> %d)\n", name, svcType, localPort, remotePort)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().String("type", "", "Service type: 'ec2' or 'k8s' (required)")
	addCmd.Flags().String("host", "", "EC2 private IP or hostname")
	addCmd.Flags().String("target", "", "K8s target: svc/<name>, deploy/<name>, or pod/<name>")
	addCmd.Flags().StringP("namespace", "n", "default", "K8s namespace")
	addCmd.Flags().String("context", "", "K8s context (default: current context)")
	addCmd.Flags().IntP("local", "l", 0, "Local port (required)")
	addCmd.Flags().IntP("remote", "r", 0, "Remote port (required)")
	addCmd.Flags().String("desc", "", "Description")
	addCmd.Flags().String("ssh-command", "", "Override SSH command for this service")
	addCmd.Flags().String("user", "", "SSH user for EC2 connection (e.g. ec2-user, ubuntu)")

	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("local")
	addCmd.MarkFlagRequired("remote")
}
