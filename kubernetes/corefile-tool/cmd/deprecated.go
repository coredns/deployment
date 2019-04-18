package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewDeprecatedCmd represents the deprecated command
func NewDeprecatedCmd() *cobra.Command {
	deprecatedCmd := &cobra.Command{
		Use:   "deprecated",
		Short: "Deprecated returns a list of deprecated plugins or directives present in the Corefile.",
		Example: `# See deprecated plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool deprecated --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")
			deprecated, err := deprecatedCorefileFromPath(from, to, corefile)
			if err != nil {
				return fmt.Errorf("error while listing deprecated plugins: %v \n", err)
			}
			for _, dep := range deprecated {
				fmt.Println(dep.ToString())
			}
			return nil
		},
	}
	deprecatedCmd.Flags().String("from", "", "Required: The version you are migrating from. ")
	deprecatedCmd.MarkFlagRequired("from")
	deprecatedCmd.Flags().String("to", "", "Required: The version you are migrating to.")
	deprecatedCmd.MarkFlagRequired("to")
	deprecatedCmd.Flags().String("corefile", "", "Required: The path where your Corefile is located.")
	deprecatedCmd.MarkFlagRequired("corefile")

	return deprecatedCmd
}

// deprecatedCorefileFromPath takes the path where the Corefile is located and returns the deprecated plugins or directives
// present in the Corefile.
func deprecatedCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string) ([]migration.Notice, error) {
	if _, err := os.Stat(corefilePath); os.IsNotExist(err) {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(corefilePath)
	if err != nil {
		return nil, err
	}
	corefileStr := string(fileBytes)
	return migration.Deprecated(fromCoreDNSVersion, toCoreDNSVersion, corefileStr)
}
