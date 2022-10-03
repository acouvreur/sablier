package cmd

import (
	"fmt"

	"github.com/acouvreur/sablier/v2/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version Sablier",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
	},
}
