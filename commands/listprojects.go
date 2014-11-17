package commands

import (
	"fmt"

	"github.com/ecornell/tfs-tool/lib"
	"github.com/spf13/cobra"
)

var cmdListProject = &cobra.Command{
	Use:   "listprojects",
	Short: "List TFS Proejects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listProjects()
	},
}

func listProjects() {

	projectsJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projectsJSON.Get("value").MustArray() {

		p := pi.(map[string]interface{})

		fmt.Println(p["name"])

	}

}
