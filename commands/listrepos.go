package commands

import (
	"fmt"
	"sort"

	"github.com/ecornell/tfs-tool/lib"
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

	projects := lib.CallTFS("/" + lib.Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	data := []projectT{}

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		name := p["name"].(string)

		pp := projectT{name: name}

		reposJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/git/" + name + "/repositories?api-version=1.0-preview.1")

		rData := []repositoryT{}

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			rData = append(rData,
				repositoryT{
					remoteURL: r["remoteUrl"].(string),
					name:      r["name"].(string),
				},
			)
		}

		sort.Sort(byNameRepositoryT(rData))
		pp.repository = rData
		data = append(data, pp)

	}

	sort.Sort(byNameProjectT(data))

	for _, project := range data {

		if lib.Flags.Color {
			color.Println("@{c}" + project.name)
		} else {
			fmt.Println(project.name)
		}

		for _, repo := range project.repository {
			if lib.Flags.Color {
				color.Println("  @g" + repo.name + " @y-> @w" + repo.remoteURL)
			} else {
				fmt.Println("  " + repo.name + " -> " + repo.remoteURL)
			}

		}
	}

}
