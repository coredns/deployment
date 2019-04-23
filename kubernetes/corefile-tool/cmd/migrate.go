package cmd

import (
	"fmt"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewMigrateCmd represents the migrate command
func NewMigrateCmd() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate your CoreDNS corefile",
		Example: `# Migrate CoreDNS from v1.4.0 to v1.5.0 and handle deprecations . 
corefile-tool migrate --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile  --deprecations true

# Migrate CoreDNS from v1.2.2 to v1.3.1 and do not handle deprecations .
corefile-tool migrate --from 1.2.2 --to 1.3.1 --corefile /path/to/Corefile  --deprecations false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")
			deprecations, _ := cmd.Flags().GetBool("deprecations")

			migrated, err := migrateCorefileFromPath(from, to, corefile, deprecations)
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
	migrateCmd.Flags().Bool("deprecations", false, "Required: Specify whether you want to handle plugin deprecations. [True | False] ")
	migrateCmd.MarkFlagRequired("deprecations")

	return migrateCmd
}

// migrateCorefileFromPath takes the path where the Corefile is located and returns the deprecated plugins or directives
// present in the Corefile.
func migrateCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string, deprecations bool) (string, error) {
	fileBytes, err := getCorefileFromPath(corefilePath)
	if err != nil {
		return "", err
	}
	corefileStr := string(fileBytes)
	return migration.Migrate(fromCoreDNSVersion, toCoreDNSVersion, corefileStr, deprecations)
}
