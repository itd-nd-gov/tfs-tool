package commands

import (
	"fmt"
	"sort"

	"github.com/ecornell/tfs-tool/lib"
	"github.com/spf13/cobra"
	"github.com/wsxiaoys/terminal/color"
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

	data := []projectT{}

	for _, pi := range projectsJSON.Get("value").MustArray() {

		p := pi.(map[string]interface{})

		pp := projectT{name: p["name"].(string)}

		data = append(data, pp)

	}

	sort.Sort(byNameProjectT(data))

	for _, project := range data {

		if lib.Flags.Color {
			color.Println("@{c}" + project.name)
		} else {
			fmt.Println(project.name)
		}

	}

}
