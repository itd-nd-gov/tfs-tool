package lib

import "strings"

func GitClone(url string) {
	execmd("git clone " + addGitAuthToRemoteURL(url))
}

func GitPull(url string) {
	execmd("git pull " + addGitAuthToRemoteURL(url))
}

func addGitAuthToRemoteURL(url string) string {
	return strings.Replace(url, "://", "://"+getUser()+":"+getPassword()+"@", -1)
}
