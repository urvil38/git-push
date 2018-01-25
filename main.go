package main

import (
	"github.com/urvil38/git-cli/git"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urvil38/git-cli/color"
	"github.com/urvil38/git-cli/encoding"
	"github.com/urvil38/git-cli/github"
	"github.com/urvil38/git-cli/questions"
	"github.com/urvil38/git-cli/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

var gitcliASCII = `
  ________.__  __                   .__  .__ 
 /  _____/|__|/  |_            ____ |  | |__|
/   \  ___|  \   __\  ______ _/ ___\|  | |  |
\    \_\  \  ||  |   /_____/ \  \___|  |_|  |
 \______  /__||__|            \___  >____/__|
        \/                        \/         
`

func init() {
	createDir()
	checkCredential()
}

func createDir() {
	err := os.MkdirAll(github.Home+"/.config/"+"git-cli", 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func checkCredential() {
	b, err := ioutil.ReadFile(github.Home + "/.config/" + "git-cli/git-cli")
	if err != nil {
		return
	}
	b = encoding.Decode(string(b))
	credentials := strings.Split(string(b), "\n")
	answer.Username = credentials[0]
	answer.Password = credentials[1]
}

var answer types.Answer

func main() {
	fmt.Println(color.Wrap(gitcliASCII, "FgYellow", "Bold"))

	err := survey.Ask(questions.ServiceName, &answer)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = survey.Ask(questions.GithubRepoInfo, &answer.Repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch service := answer.ServiceName; service {
	case "Github":
		err := github.Init(&answer)
		if err != nil {
			fmt.Println(color.Wrap("=> "+err.Error(), "FgRed", "CrossedOut"))
			os.Exit(0)
		}
		err = github.CreateRepo(&answer)
		if err != nil {
			fmt.Println(color.Wrap("=> "+err.Error(), "FgRed", "CrossedOut"))
			os.Exit(0)
		}
		err = git.CreateGitIgnoreFile()
	}

}
