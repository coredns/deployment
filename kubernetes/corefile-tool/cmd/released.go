package cmd

import (
	"fmt"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// NewReleasedCmd represents the released command
func NewReleasedCmd() *cobra.Command {
	releasedCmd := &cobra.Command{
		Use:   "released",
		Short: "Determines whether your Docker Image ID of a CoreDNS release is valid or not",
		Run: func(cmd *cobra.Command, args []string) {
			image, _ := cmd.Flags().GetString("dockerImageID")
			result := migration.Released(image)

			if result {
				fmt.Println("The docker image ID is valid")
			} else {
				fmt.Println("The docker image ID is invalid")
			}
		},
	}

	releasedCmd.Flags().String("dockerImageID", "", "Required: The docker image ID you want to check. ")
	releasedCmd.MarkFlagRequired("dockerImageID")

	return releasedCmd
}
