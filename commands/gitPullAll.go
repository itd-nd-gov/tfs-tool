package commands

import (
	"fmt"
	"os"

	"github.com/ecornell/tfs-tool/lib"
	"github.com/ecornell/tfs-tool/utils"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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
		jww.ERROR.Println("Output directory required")
		return
	}

	var err error
	err = os.MkdirAll(lib.Flags.DestinationDir, 0777)
	if err != nil {
		utils.CheckErr(err)
		return
	}

	err = os.Chdir(lib.Flags.DestinationDir)
	if err != nil {
		utils.CheckErr(err)
		return
	}

	baseDir, _ := os.Getwd()

	projects := lib.CallTFS("/" + lib.Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	// pull all
	var projectNames []string
	var repoNames []string

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		projectNames = append(projectNames, projectName)

		reposJSON := lib.CallTFS("/" + lib.Config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			remoteURL := r["remoteUrl"].(string)
			name := r["name"].(string)

			repoNames = append(repoNames, name)

			os.Chdir(baseDir)
			os.Mkdir(projectName, 0777)
			os.Chdir(projectName)

			err := os.Chdir(name)
			if err != nil {
				fmt.Println("Cloning - " + projectName + " :: " + name)
				lib.GitClone(remoteURL)
			} else {
				os.Chdir(name)
				fmt.Println("Pulling - " + projectName + " :: " + name)
				lib.GitPull(remoteURL)
			}

		}

		// cleanup deleted repos
		cleanDir(baseDir+"/"+projectName, repoNames)
	}

	// cleanup deleted team projects
	cleanDir(baseDir, projectNames)

}

func cleanDir(baseDir string, validSubDirs []string) {
	// cleanup deleted team projects
	dir, err := os.Open(baseDir)
	if err != nil {
		utils.CheckErr(err)
		return
	}
	// checkErr(err)
	defer dir.Close()
	fi, err := dir.Stat()
	if err != nil {
		utils.CheckErr(err)
		return
	}
	var dirnames []string
	if fi.IsDir() {
		fis, err := dir.Readdir(-1) // -1 means return all the FileInfos
		if err != nil {
			utils.CheckErr(err)
			return
		}
		for _, fileinfo := range fis {
			if fileinfo.IsDir() {
				dirnames = append(dirnames, fileinfo.Name())
			}
		}
	}
	os.Chdir(baseDir)
	for _, d := range dirnames {
		if !stringInSlice(d, validSubDirs) {
			os.Remove(d)
			jww.INFO.Println("Removed - " + baseDir + "/" + d)
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
