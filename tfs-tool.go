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

const Version = "0.1"

type ConfigT struct {
	BaseURL    string
	User       string
	Password   string
	Collection string
}

var Config ConfigT

type FlagsT struct {
	UserID         string
	Password       string
	DestinationDir string
	Verbose        bool
	Color          bool
}

var Flags FlagsT

var key = []byte("caskd92h3jfld0u3jlaafsd08jz2cv3a")

func main() {

	_, err := toml.DecodeFile("tfs-tool.cfg", &Config)
	check(err)

	if !strings.HasPrefix(Config.Password, "~~") {
		ciphertext, err := encrypt(key, []byte(Config.Password))
		check(err)

		Config.Password = "~~" + hex.EncodeToString(ciphertext)

		var ConfigBuffer bytes.Buffer
		toml.NewEncoder(&ConfigBuffer).Encode(Config)

		err = ioutil.WriteFile("tfs-tool.cfg", ConfigBuffer.Bytes(), 0644)
		check(err)

	} else {
		encPassword, _ := hex.DecodeString(Config.Password[2:])
		result, err := decrypt(key, encPassword)
		check(err)
		Config.Password = string(result)
	}

	//

	var CmdListProject = &cobra.Command{
		Use:   "listprojects",
		Short: "List TFS Proejects",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listProjects()
		},
	}

	var CmdListRepos = &cobra.Command{
		Use:   "listrepos",
		Short: "List TFS repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listRepos()
		},
	}

	var CmdPullAll = &cobra.Command{
		Use:   "pullall",
		Short: "Pull all TFS repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			pullAll()
		},
	}
	CmdPullAll.Flags().StringVarP(&Flags.DestinationDir, "dest", "d", "", "Output directory to store repositories")

	var CmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of tfs-tool",
		Long:  `All software has versions. This is tfs-tool's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}

	var RootCmd = &cobra.Command{Use: "tfs-tool"}
	RootCmd.AddCommand(CmdListProject, CmdListRepos, CmdPullAll, CmdVersion)

	RootCmd.PersistentFlags().BoolVarP(&Flags.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&Flags.Color, "color", "", false, "colorize output")

	RootCmd.PersistentFlags().StringVarP(&Flags.UserID, "user", "", "", "TFS User ID")

	RootCmd.PersistentFlags().StringVarP(&Flags.Password, "password", "", "", "TFS Password")

	RootCmd.Execute()

}

func listProjects() {

	projectsJSON := callTFS("/" + Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projectsJSON.Get("value").MustArray() {

		p := pi.(map[string]interface{})

		fmt.Println(p["name"])

	}

}

func listRepos() {

	projects := callTFS("/" + Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		reposJSON := callTFS("/" + Config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

		if Flags.Color {
			color.Println("@{c}" + projectName)
		} else {
			fmt.Println(projectName)
		}

		for _, ri := range reposJSON.Get("value").MustArray() {
			r, _ := ri.(map[string]interface{})

			remoteURL := r["remoteUrl"].(string)
			name := r["name"].(string)

			if Flags.Color {
				color.Println("  @g" + name + " @y-> @w" + remoteURL)
			} else {
				fmt.Println("  " + name + " -> " + remoteURL)
			}
		}

	}

}

func pullAll() {

	if Flags.DestinationDir == "" {
		fmt.Println("ERROR: Output directory required")
		return
	}

	os.MkdirAll(Flags.DestinationDir, 0777)
	os.Chdir(Flags.DestinationDir)
	baseDir, _ := os.Getwd()

	projects := callTFS("/" + Config.Collection + "/_apis/projects?api-version=1.0-preview.2")

	for _, pi := range projects.Get("value").MustArray() {
		p, _ := pi.(map[string]interface{})

		projectName := p["name"].(string)

		reposJSON := callTFS("/" + Config.Collection + "/_apis/git/" + projectName + "/repositories?api-version=1.0-preview.1")

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
				exeCmd("git clone " + remoteURLAuth)
			} else {
				os.Chdir(name)
				fmt.Println("Pulling - " + name)
				exeCmd("git pull " + remoteURLAuth)
			}

		}
	}
}

func addGitAuthToRemoteURL(url string) string {
	return strings.Replace(url, "://", "://"+getUser()+":"+getPassword()+"@", -1)
}

func callTFS(cmd string) *simplejson.Json {
	var curl = "curl -s --ntlm -k --negotiate -u " + getUser() + ":" + getPassword() + " "

	output := exeCmd(curl + Config.BaseURL + cmd)

	js, _ := simplejson.NewJson([]byte(output))

	return js
}

func exeCmd(cmd string) string {

	if Flags.Verbose {
		fmt.Println("command is ", cmd)
	}

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	if Flags.Verbose {
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
	if Flags.UserID != "" {
		return Flags.UserID
	} else {
		return Config.User
	}
}

func getPassword() string {
	if Flags.Password != "" {
		return Flags.Password
	} else {
		return Config.Password
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
