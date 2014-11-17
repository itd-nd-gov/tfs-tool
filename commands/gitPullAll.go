package commands

import (
	"fmt"
	"os"

	"github.com/ecornell/tfs-tool/lib"
	"github.com/spf13/cobra"
)

var cmdGitPullAll = &cobra.Command{
	Use:   "gitpullall",
	Short: "Pull all TFS Git repositories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pullAll()
	},
}

func init() {
	cmdGitPullAll.Flags().StringVarP(&lib.Flags.DestinationDir, "dest", "d", "", "Output directory to store repositories")
}

func pullAll() {

	if lib.Flags.DestinationDir == "" {
		fmt.Println("ERROR: Output directory required")
		return
	}

	os.MkdirAll(lib.Flags.DestinationDir, 0777)
	os.Chdir(lib.Flags.DestinationDir)
	baseDir, _ := os.Getwd()

	projects := lib.CallTFS("/" + lib.Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		reposJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			remoteURL := r["remoteUrl"].(string)
			name := r["name"].(string)

			os.Chdir(baseDir)
			os.Mkdir(projectName, 0777)
			os.Chdir(projectName)

			err := os.Chdir(name)
			if err != nil {
				fmt.Println("Cloning - " + name)
				lib.GitClone(remoteURL)
			} else {
				os.Chdir(name)
				fmt.Println("Pulling - " + name)
				lib.GitPull(remoteURL)
			}

		}
	}
}
