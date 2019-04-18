package cmd

import (
	"fmt"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command for the corefile-tool.
var rootCmd = &cobra.Command{
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

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
