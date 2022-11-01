package cmd

import (
	"cli/actions"
	"cli/env"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var clusterName string

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the cluster",
		Long:  `Destroys the cluster with a given name.`,

		Run: func(cmd *cobra.Command, args []string) {
			err := actions.Destroy(clusterName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "specify the cluster to be used")
	destroyCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")
	destroyCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	// Auto complete cluster names of active clusters for flag 'cluster'.
	destroyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		clusters, err := actions.ReadClustersInfo()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var names []string
		for _, c := range clusters {
			if c.Active() {
				names = append(names, c.Name)
			}
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	})
}
