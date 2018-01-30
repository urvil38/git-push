package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/bitbucket"
	"github.com/urvil38/git-push/git"
	"github.com/urvil38/git-push/github"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/urvil38/git-push/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	colorRed = color.New(color.FgRed, color.Bold)
	colorYellow = color.New(color.FgYellow, color.Bold)
	home = os.Getenv("HOME")
	if home == "" {
		fmt.Println(help)
		os.Exit(0)
	}
	userConfigFile = home + separator + ".config" + separator + "git-push" + separator + "userInfo"
	configFolder = home + separator + ".config" + separator + "git-push"
	createDir()
	checkUserInfo()
	remoteExists, _ = utils.CheckRemoteRepo()
	if remoteExists {
		colorRed.Println("Sorry, this tool will not help you because working repository is already on github or bitbucket!")
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
		colorRed.Println("=> " + err.Error())
		os.Exit(0)
	}
}

var (
	serviceName    string
	repo           types.Repo
	basicUserInfo  types.BasicUserInfo
	remoteExists   bool
	err            error
	home           string
	userConfigFile string
	configFolder   string
	colorRed       *color.Color
	colorYellow    *color.Color
)

const (
	banner = `
  ________ .__   __              __________                .__     
 /  _____/ |__|_/  |_            \______   \ __ __   ______|  |__  
/   \  ___ |  |\   __\   ______   |     ___/|  |  \ /  ___/|  |  \ 
\    \_\  \|  | |  |    /_____/   |    |    |  |  / \___ \ |   Y  \
 \______  /|__| |__|              |____|    |____/ /____  >|___|  /
     	\/                                              \/      \/ 
`
	separator = string(filepath.Separator)
	help      = `
---------------x configure x----------------

For linux and macos:

	export $HOME=/path/to/home/where/git-push/can/store/credentials

For windows:

	you must set the HOME environment variable to your chosen path(I suggest c:\git-push)
	
	Under Windows, you may set environment variables through the "Environment Variables" 
	button on the "Advanced" tab of the "System" control panel. Some versions of Windows 
	provide this control panel through the "Advanced System Settings" option inside 
	the "System" control panel. 	
`
)

func main() {
	colorYellow.Println(banner)

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

	err = survey.Ask(questions.ServiceName, &serviceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !remoteExists {
		err = survey.Ask(questions.GithubRepoInfo, &repo)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch serviceName {
	case "Github":
		err := github.Init()
		checkerror(err)

		err = github.CreateRepo(repo)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(github.GitURL, github.GithubUser, basicUserInfo)
		checkerror(err)
	case "Bitbucket":
		err := bitbucket.Init()
		checkerror(err)

		err = bitbucket.CreateRepo(repo)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(bitbucket.BitbuckerURL, bitbucket.BitbucketUser, basicUserInfo)
		checkerror(err)
	}
}
