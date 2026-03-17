package cmd

import (
	"fmt"

	"github.com/keimoon/pf/config"
	"github.com/spf13/cobra"
)

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "View or set default configuration",
	Long: `View or set defaults. Examples:

  pf defaults                             # show current defaults
  pf defaults --ssh-command ssh-nohost    # set default SSH command
  pf defaults --user ec2-user             # set default SSH user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configPath()
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		sshCmd, _ := cmd.Flags().GetString("ssh-command")
		sshUser, _ := cmd.Flags().GetString("user")

		if sshCmd == "" && sshUser == "" {
			// Show mode
			fmt.Printf("ssh_command: %s\n", cfg.Defaults.SSHCommand)
			if cfg.Defaults.SSHUser != "" {
				fmt.Printf("ssh_user: %s\n", cfg.Defaults.SSHUser)
			}
			return nil
		}

		// Set mode
		if sshCmd != "" {
			cfg.Defaults.SSHCommand = sshCmd
			fmt.Printf("Set ssh_command to %q\n", sshCmd)
		}
		if sshUser != "" {
			cfg.Defaults.SSHUser = sshUser
			fmt.Printf("Set ssh_user to %q\n", sshUser)
		}
		if err := config.Save(path, cfg); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(defaultsCmd)
	defaultsCmd.Flags().String("ssh-command", "", "Default SSH command (e.g. ssh, ssh-nohost)")
	defaultsCmd.Flags().String("user", "", "Default SSH user for EC2 connections (e.g. ec2-user, ubuntu)")
}
