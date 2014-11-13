package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/bitly/go-simplejson"
	"github.com/spf13/cobra"
	"github.com/wsxiaoys/terminal/color"
)

const version = "0.1"

type configT struct {
	BaseURL    string
	User       string
	Password   string
	Collection string
}

var config configT

type flagsT struct {
	UserID         string
	Password       string
	DestinationDir string
	Verbose        bool
	Color          bool
}

var flags flagsT

var key = []byte("caskd92h3jfld0u3jlaafsd08jz2cv3a")

func main() {

	_, err := toml.DecodeFile("tfs-tool.cfg", &config)
	check(err)

	if !strings.HasPrefix(config.Password, "~~") {
		ciphertext, err := encrypt(key, []byte(config.Password))
		check(err)

		config.Password = "~~" + hex.EncodeToString(ciphertext)

		var configBuffer bytes.Buffer
		toml.NewEncoder(&configBuffer).Encode(config)

		err = ioutil.WriteFile("tfs-tool.cfg", configBuffer.Bytes(), 0644)
		check(err)

	} else {
		encPassword, _ := hex.DecodeString(config.Password[2:])
		result, err := decrypt(key, encPassword)
		check(err)
		config.Password = string(result)
	}

	//

	var cmdListProject = &cobra.Command{
		Use:   "listprojects",
		Short: "List TFS Proejects",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listProjects()
		},
	}

	var cmdListRepos = &cobra.Command{
		Use:   "listrepos",
		Short: "List TFS repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listRepos()
		},
	}

	var cmdPullAll = &cobra.Command{
		Use:   "pullall",
		Short: "Pull all TFS repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			pullAll()
		},
	}
	cmdPullAll.Flags().StringVarP(&flags.DestinationDir, "dest", "d", "", "Output directory to store repositories")

	var cmdversion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of tfs-tool",
		Long:  `All software has versions. This is tfs-tool's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var cmdRoot = &cobra.Command{Use: "tfs-tool"}
	cmdRoot.AddCommand(cmdListProject, cmdListRepos, cmdPullAll, cmdversion)

	cmdRoot.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	cmdRoot.PersistentFlags().BoolVarP(&flags.Color, "color", "", false, "colorize output")
	cmdRoot.PersistentFlags().StringVarP(&flags.UserID, "user", "", "", "TFS User ID")
	cmdRoot.PersistentFlags().StringVarP(&flags.Password, "password", "", "", "TFS Password")

	cmdRoot.Execute()

}

func listProjects() {

	projectsJSON := callTFS("/" + config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projectsJSON.Get("value").MustArray() {

		p := pi.(map[string]interface{})

		fmt.Println(p["name"])

	}

}

func listRepos() {

	projects := callTFS("/" + config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		reposJSON := callTFS("/" + config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

		if flags.Color {
			color.Println("@{c}" + projectName)
		} else {
			fmt.Println(projectName)
		}

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			remoteURL := r["remoteUrl"].(string)
			name := r["name"].(string)

			if flags.Color {
				color.Println("  @g" + name + " @y-> @w" + remoteURL)
			} else {
				fmt.Println("  " + name + " -> " + remoteURL)
			}
		}

	}

}

func pullAll() {

	if flags.DestinationDir == "" {
		fmt.Println("ERROR: Output directory required")
		return
	}

	os.MkdirAll(flags.DestinationDir, 0777)
	os.Chdir(flags.DestinationDir)
	baseDir, _ := os.Getwd()

	projects := callTFS("/" + config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		reposJSON := callTFS("/" + config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			remoteURL := r["remoteUrl"].(string)
			name := r["name"].(string)

			remoteURLAuth := addGitAuthToRemoteURL(remoteURL)

			os.Chdir(baseDir)
			os.Mkdir(projectName, 0777)
			os.Chdir(projectName)

			err := os.Chdir(name)
			if err != nil {
				fmt.Println("Cloning - " + name)
				execmd("git clone " + remoteURLAuth)
			} else {
				os.Chdir(name)
				fmt.Println("Pulling - " + name)
				execmd("git pull " + remoteURLAuth)
			}

		}
	}
}

func addGitAuthToRemoteURL(url string) string {
	return strings.Replace(url, "://", "://"+getUser()+":"+getPassword()+"@", -1)
}

func callTFS(cmd string) *simplejson.Json {
	var curl = "curl -s --ntlm -k --negotiate -u " + getUser() + ":" + getPassword() + " "

	output := execmd(curl + config.BaseURL + cmd)

	js, _ := simplejson.NewJson([]byte(output))

	return js
}

func execmd(cmd string) string {

	if flags.Verbose {
		fmt.Println("command is ", cmd)
	}

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	if flags.Verbose {
		fmt.Printf("%s\n", out)
	}

	return string(out)
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getUser() string {
	if flags.UserID != "" {
		return flags.UserID
	}
	return config.User
}

func getPassword() string {
	if flags.Password != "" {
		return flags.Password
	}
	return config.Password
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
