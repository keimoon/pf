package cmd

import (
	"fmt"
	"os"

	"github.com/keimoon/pf/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pf",
	Short: "Port-forward manager for k8s services and EC2 instances",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.config/pf/services.yaml)")
}

func configPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return config.DefaultPath()
}
