package cmd

import (
	"fmt"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewDowngradeCmd represents the downgrade command
func NewDowngradeCmd() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "downgrade",
		Short: "Downgrade your CoreDNS corefile to a previous version",
		Example: `# Downgrade CoreDNS from v1.5.0 to v1.4.0. 
corefile-tool downgrade --from 1.5.0 --to 1.4.0 --corefile /path/to/Corefile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")

			migrated, err := downgradeCorefileFromPath(from, to, corefile)
			if err != nil {
				return fmt.Errorf("error while migration: %v \n", err)
			}
			fmt.Println(migrated)
			return nil
		},
	}
	migrateCmd.Flags().String("from", "", "Required: The version you are migrating from. ")
	migrateCmd.MarkFlagRequired("from")
	migrateCmd.Flags().String("to", "", "Required: The version you are migrating to.")
	migrateCmd.MarkFlagRequired("to")
	migrateCmd.Flags().String("corefile", "", "Required: The path where your Corefile is located.")
	migrateCmd.MarkFlagRequired("corefile")
	migrateCmd.Flags().Bool("deprecations", false, "Specify whether you want to handle plugin deprecations. [True | False] ")

	return migrateCmd
}

// downgradeCorefileFromPath takes the path where the Corefile is located and downgrades the Corefile to the
// desrired version.
func downgradeCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string) (string, error) {
	fileBytes, err := getCorefileFromPath(corefilePath)
	if err != nil {
		return "", err
	}
	corefileStr := string(fileBytes)
	return migration.MigrateDown(fromCoreDNSVersion, toCoreDNSVersion, corefileStr)
}
