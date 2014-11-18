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

	projects := []project{}

	for _, pi := range projectsJSON.Get("value").MustArray() {

		p := pi.(map[string]interface{})

		project := project{name: p["name"].(string)}

		projects = append(projects, project)

	}

	sort.Sort(projectsByName{projects})

	for _, project := range projects {

		if lib.Flags.Color {
			color.Println("@{c}" + project.name)
		} else {
			fmt.Println(project.name)
		}

	}

}
