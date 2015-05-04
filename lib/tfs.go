package lib

import (
	"os/exec"
	"strings"

	"github.com/bitly/go-simplejson"
	jww "github.com/spf13/jwalterweatherman"
)

func CallTFS(cmd string) *simplejson.Json {
	var curl = "curl -s --ntlm -k --negotiate -u " + getUser() + ":" + getPassword() + " "

	output := execmd(curl + Config.BaseURL + cmd)

	js, _ := simplejson.NewJson([]byte(output))

	return js
}

func execmd(cmd string) string {

	jww.TRACE.Println("command is ", cmd)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		jww.WARN.Printf("%s\n", err)
	}

	jww.TRACE.Printf("%s\n", out)

	return string(out)
}
