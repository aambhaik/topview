package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of topview CLI",
	Long:  `All software has versions. This is topview CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mashery Topology Viewer v0.1")
	},
}
