package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/keimoon/pf/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all registered services",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configPath()
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		if len(cfg.Services) == 0 {
			fmt.Println("No services registered. Use 'pf add' to add one.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tTARGET\tLOCAL\tREMOTE\tDESCRIPTION")

		names := make([]string, 0, len(cfg.Services))
		for name := range cfg.Services {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			svc := cfg.Services[name]
			target := svc.Host
			if svc.Type == "k8s" {
				target = svc.Target
				if svc.Namespace != "" {
					target = svc.Namespace + "/" + svc.Target
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\n",
				name, svc.Type, target, svc.LocalPort, svc.RemotePort, svc.Description)
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
