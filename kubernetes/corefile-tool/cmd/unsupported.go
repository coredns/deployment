package cmd

import (
	"fmt"
	"github.com/coredns/deployment/kubernetes/migration"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// NewUnsupportedCmd represents the unsupported command
func NewUnsupportedCmd() *cobra.Command {
	unsupportedCmd := &cobra.Command{
		Use:   "unsupported",
		Short: "Unsupported returns a list of plugins that are not recognized/supported by the migration tool (but may still be valid in CoreDNS).",
		Example: `# See unsupported plugins CoreDNS from v1.4.0 to v1.5.0. 
corefile-tool unsupported --from 1.4.0 --to 1.5.0 --corefile /path/to/Corefile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			corefile, _ := cmd.Flags().GetString("corefile")
			unsupported, err := unsupportedCorefileFromPath(from, to, corefile)
			if err != nil {
				return fmt.Errorf("error while listing deprecated plugins: %v \n", err)
			}
			for _, unsup := range unsupported {
				fmt.Println(unsup.ToString())
			}
			return nil
		},
	}

	unsupportedCmd.Flags().String("from", "", "Required: The version you are migrating from. ")
	unsupportedCmd.MarkFlagRequired("from")
	unsupportedCmd.Flags().String("to", "", "Required: The version you are migrating to.")
	unsupportedCmd.MarkFlagRequired("to")
	unsupportedCmd.Flags().String("corefile", "", "Required: The path where your Corefile is located.")
	unsupportedCmd.MarkFlagRequired("corefile")

	return unsupportedCmd
}

// unsupportedCorefileFromPath takes the path where the Corefile is located and returns a list of  plugins
// that have been removed.
func unsupportedCorefileFromPath(fromCoreDNSVersion, toCoreDNSVersion, corefilePath string) ([]migration.Notice, error) {
	if _, err := os.Stat(corefilePath); os.IsNotExist(err) {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(corefilePath)
	if err != nil {
		return nil, err
	}
	corefileStr := string(fileBytes)
	return migration.Unsupported(fromCoreDNSVersion, toCoreDNSVersion, corefileStr)
}
