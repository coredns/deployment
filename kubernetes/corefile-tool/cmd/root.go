package cmd

import (
	"fmt"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// CorefileTool represents the base command for the corefile-tool.
func CorefileTool() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "corefile-tool",
		Short: "A brief description of your application",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ CoreDNS Migration Tool                                   │
			    │ Easily Migrate your Corefile                             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/coredns/deployment/issues             │
			    └──────────────────────────────────────────────────────────┘

		`),
	}
	rootCmd.AddCommand(NewRemovedCmd())
	rootCmd.AddCommand(NewMigrateCmd())
	rootCmd.AddCommand(NewDefaultCmd())
	rootCmd.AddCommand(NewDeprecatedCmd())
	rootCmd.AddCommand(NewUnsupportedCmd())
	rootCmd.AddCommand(NewValidVersionsCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := CorefileTool().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
