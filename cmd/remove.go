package cmd

import (
	"fmt"

	"github.com/keimoon/pf/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a registered service",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := configPath()

		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		if err := cfg.Remove(name); err != nil {
			return err
		}

		if err := config.Save(path, cfg); err != nil {
			return err
		}

		fmt.Printf("Removed %q\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
