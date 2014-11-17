package lib

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bitly/go-simplejson"
)

func CallTFS(cmd string) *simplejson.Json {
	var curl = "curl -s --ntlm -k --negotiate -u " + getUser() + ":" + getPassword() + " "

	output := execmd(curl + Config.BaseURL + cmd)

	js, _ := simplejson.NewJson([]byte(output))

	return js
}

func execmd(cmd string) string {

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
