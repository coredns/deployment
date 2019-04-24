package cmd

import (
	"fmt"
	"strings"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewValidVersionsCmd represents the validversions command
func NewValidVersionsCmd() *cobra.Command {
	validversionsCmd := &cobra.Command{
		Use:   "validversions",
		Short: "Shows valid versions of CoreDNS",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The following are valid CoreDNS versions:")
			fmt.Println(strings.Join(migration.ValidVersions(), ", "))
		},
	}
	return validversionsCmd
}
