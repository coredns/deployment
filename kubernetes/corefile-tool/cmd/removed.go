package cmd

import (
	"fmt"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewRemovedCmd represents the removed command
func NewRemovedCmd() *cobra.Command {
	removedCmd := &cobra.Command{
		Use:   "removed",
		Short: "Removed returns a list of removed plugins or directives present in the Corefile.",
		Example: `# See removed plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool removed --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")
			removed, err := removedCorefileFromPath(from, to, corefile)
			if err != nil {
				return fmt.Errorf("error while listing deprecated plugins: %v \n", err)
			}
			for _, rem := range removed {
				fmt.Println(rem.ToString())
			}
			return nil
		},

	}
	removedCmd.Flags().String("from", "", "Required: The version you are migrating from. ")
	removedCmd.MarkFlagRequired("from")
	removedCmd.Flags().String("to", "", "Required: The version you are migrating to.")
	removedCmd.MarkFlagRequired("to")
	removedCmd.Flags().String("corefile", "", "Required: The path where your Corefile is located.")
	removedCmd.MarkFlagRequired("corefile")

	return removedCmd
}

// removedCorefileFromPath takes the path where the Corefile is located and returns the plugins or directives
// that have been removed.
func removedCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string) ([]migration.Notice, error) {
	fileBytes, err := getCorefileFromPath(corefilePath)
	if err != nil {
		return nil, err
	}
	corefileStr := string(fileBytes)
	return migration.Removed(fromCoreDNSVersion, toCoreDNSVersion, corefileStr)
}
