package commands

import (
	"github.com/ecornell/tfs-tool/lib"
	"github.com/ecornell/tfs-tool/utils"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var CmdRoot = &cobra.Command{Use: "tfs-tool"}

func Execute() {
	AddCommands()
	utils.StopOnErr(CmdRoot.Execute())
}

func AddCommands() {
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

	jww.SetLogFile("tfs-tool.log")

	if lib.Flags.Verbose {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelInfo)
	} else {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelError)
	}

}
