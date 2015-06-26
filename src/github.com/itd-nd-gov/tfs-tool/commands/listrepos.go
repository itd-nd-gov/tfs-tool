package commands

import (
	"fmt"
	"sort"

	"github.com/itd-nd-gov/tfs-tool/lib"
	"github.com/spf13/cobra"
	"github.com/wsxiaoys/terminal/color"
)

var cmdListRepos = &cobra.Command{
	Use:   "listrepos",
	Short: "List TFS repositories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listRepos()
	},
}

func listRepos() {

	projectsJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	projects := []project{}

	for _, pi := range projectsJSON.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		name := p["name"].(string)

		project := project{name: name}

		reposJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/git/" + name + "/repositories?api-version=1.0-preview.1")

		repositories := []repository{}

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			repositories = append(repositories,
				repository{
					remoteURL: r["remoteUrl"].(string),
					name:      r["name"].(string),
				},
			)
		}

		sort.Sort(repositoriesByName{repositories})
		project.repositories = repositories
		projects = append(projects, project)

	}

	sort.Sort(projectsByName{projects})

	for _, project := range projects {

		if lib.Flags.Color {
			color.Println("@{c}" + project.name)
		} else {
			fmt.Println(project.name)
		}

		for _, repo := range project.repositories {
			if lib.Flags.Color {
				color.Println("  @g" + repo.name + " @y-> @w" + repo.remoteURL)
			} else {
				fmt.Println("  " + repo.name + " -> " + repo.remoteURL)
			}

		}
	}

}
