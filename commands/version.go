package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.1"

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of tfs-tool",
	Long:  `All software has versions. This is tfs-tool's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
