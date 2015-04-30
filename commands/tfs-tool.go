package commands

import (
	"github.com/ecornell/tfs-tool/lib"
	"github.com/ecornell/tfs-tool/utils"
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{Use: "tfs-tool"}

func Execute() {
	AddCommands()
	utils.StopOnErr(CmdRoot.Execute())
}

func AddCommands() {
	// CmdRoot.AddCommand(cmdListProject, cmdListRepos, cmdPullAll, cmdversion)
	CmdRoot.AddCommand(cmdListRepos)
	CmdRoot.AddCommand(cmdListProject)
	CmdRoot.AddCommand(cmdGitPullAll)
	CmdRoot.AddCommand(cmdVersion)
}

func init() {

	lib.LoadConfig()

	CmdRoot.PersistentFlags().BoolVarP(&lib.Flags.Verbose, "verbose", "v", false, "verbose output")
	CmdRoot.PersistentFlags().BoolVarP(&lib.Flags.Color, "color", "", false, "colorize output")
	CmdRoot.PersistentFlags().StringVarP(&lib.Flags.UserID, "user", "", "", "TFS User ID")
	CmdRoot.PersistentFlags().StringVarP(&lib.Flags.Password, "password", "", "", "TFS Password")

}
