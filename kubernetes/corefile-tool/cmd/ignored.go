package cmd

import (
	"fmt"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewIgnoredCmd represents the ignored command
func NewIgnoredCmd() *cobra.Command {
	removedCmd := &cobra.Command{
		Use:   "ignored",
		Short: "Ignored returns a list of ignored plugins or directives present in the Corefile.",
		Example: `# See removed plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool ignored --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")
			removed, err := ignoredCorefileFromPath(from, to, corefile)
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

// ignoredCorefileFromPath takes the path where the Corefile is located and returns the plugins or directives
// that have been ignored.
func ignoredCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string) ([]migration.Notice, error) {
	fileBytes, err := getCorefileFromPath(corefilePath)
	if err != nil {
		return nil, err
	}
	corefileStr := string(fileBytes)
	return migration.Ignored(fromCoreDNSVersion, toCoreDNSVersion, corefileStr)
}
