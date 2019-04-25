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
		Short: "Determines whether your Docker Image SHA of a CoreDNS release is valid or not",
		Run: func(cmd *cobra.Command, args []string) {
			image, _ := cmd.Flags().GetString("dockerImageSHA")
			result := migration.Released(image)

			if result {
				fmt.Println("The docker image SHA is valid")
			} else {
				fmt.Println("The docker image SHA is invalid")
			}
		},
	}

	releasedCmd.Flags().String("dockerImageSHA", "", "Required: The docker image SHA you want to check. ")
	releasedCmd.MarkFlagRequired("dockerImageSHA")

	return releasedCmd
}
