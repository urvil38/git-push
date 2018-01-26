package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urvil38/git-push/utils"
	"github.com/urvil38/git-push/color"
	"github.com/urvil38/git-push/git"
	"github.com/urvil38/git-push/github"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

var gitcliASCII = `
  ________ .__   __              __________                .__     
 /  _____/ |__|_/  |_            \______   \ __ __   ______|  |__  
/   \  ___ |  |\   __\   ______   |     ___/|  |  \ /  ___/|  |  \ 
\    \_\  \|  | |  |    /_____/   |    |    |  |  / \___ \ |   Y  \
 \______  /|__| |__|              |____|    |____/ /____  >|___|  /
     	\/                                              \/      \/ 
`

func init() {
	home = os.Getenv("HOME")
	userConfigFile = home + separator + ".config" + separator + "git-push" + separator + "userInfo"
	configFolder = home + separator + ".config" + separator + "git-push"
	createDir()
	checkUserInfo()
	remoteExists, _ = utils.CheckRemoteRepo()
	if remoteExists {
		fmt.Println(color.Wrap("Sorry, this tool will not help you because working repository is already on github or bitbucket!","FgRed","Bold"))
		os.Exit(0)
	}
}

func checkUserInfo() {
	b, err := ioutil.ReadFile(userConfigFile)
	if err != nil {
		return
	}
	userInfo := strings.Split(string(b), "\n")
	basicUserInfo.Name = userInfo[0]
	basicUserInfo.Email = userInfo[1]
}

func createDir() {
	err := os.MkdirAll(configFolder, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func checkerror(err error) {
	if err != nil {
		fmt.Println(color.Wrap("=> "+err.Error(), "FgRed", "CrossedOut"))
		os.Exit(0)
	}
}

var (
	answer         types.Answer
	basicUserInfo  types.BasicUserInfo
	remoteExists   bool
	err            error
	home           string
	userConfigFile string
	configFolder   string
)

const (
	separator = string(filepath.Separator)
)

func main() {
	fmt.Println(color.Wrap(gitcliASCII, "FgYellow", "Bold"))

	if basicUserInfo.Email == "" || basicUserInfo.Name == "" {
		err := survey.Ask(questions.UserInfo, &basicUserInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile(userConfigFile, []byte(basicUserInfo.Name+"\n"+basicUserInfo.Email), 0555)
		if err != nil {
			return
		}
	}

	err = survey.Ask(questions.ServiceName, &answer)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !remoteExists {
		err = survey.Ask(questions.GithubRepoInfo, &answer.Repo)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch service := answer.ServiceName; service {
	case "Github":
		err := github.Init(answer)
		checkerror(err)

		err = github.CreateRepo(answer)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(github.GitURL, github.GithubUser, basicUserInfo)
		checkerror(err)
	}
}
