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
	"github.com/urvil38/git-push/gitlab"
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

	userConfigFile = filepath.Join(home, ".config", "git-push", "userInfo")
	configFolder = filepath.Join(home, ".config", "git-push")

	createDir()
	checkUserInfo()

	remoteExists, _ = utils.CheckRemoteRepo()
	if remoteExists {
		colorRed.Println("Sorry, this tool will not help you because working repository is already on github or bitbucket or gitlab!")
		colorRed.Print("ℹ  You can use ")
		colorYellow.Print("$ git push origin master")
		colorRed.Println(" to push changes.")
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
	basicUserInfo  types.BasicUserInfo
	remoteExists   bool
	err            error
	home           string
	userConfigFile string
	configFolder   string
	version        string
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

 # Author   :  Urvil Patel
 # Version  :  %s	
 # Twitter  :  @UrvilPatel12
 # Github   :  https://github.com/urvil38
`
	help = `
***************************************| configure |******************************************

For linux and macos:
-------------------
	
	export HOME=/path/to/home/where/git-push/can/store/credentials

For windows:
-----------
	
	You must set the HOME environment variable to your chosen path(I suggest c:\git-push)

	There are two ways to doing this:
	---------------------------------

	1. Using Command Prompt you can set this environment variable by following command:
        
        set HOME="c:\git-push" 
	
	2. Under Windows, you may set environment variables through the "Environment Variables" 
	button on the "Advanced" tab of the "System" control panel. Some versions of Windows 
	provide this control panel through the "Advanced System Settings" option inside 
	the "System" control panel. 	
`
)

func main() {

	colorYellow.Println(fmt.Sprintf(banner, version))

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

	var serviceName string
	err = survey.Ask(questions.ServiceName, &serviceName)
	if err != nil {
		fmt.Println(err)
		return
	}

	var repo types.Repo
	if !remoteExists {
		err = survey.Ask(questions.GithubRepoInfo, &repo)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch serviceName {
	case "GitHub":
		err := github.Init()
		checkerror(err)

		err = github.CreateRepo(repo)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(github.GitURL, github.GithubUser, basicUserInfo)
		checkerror(err)
	case "BitBucket":
		err := bitbucket.Init()
		checkerror(err)

		err = bitbucket.CreateRepo(repo)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(bitbucket.BitbuckerURL, bitbucket.BitbucketUser, basicUserInfo)
		checkerror(err)
	case "GitLab":
		err := gitlab.Init()
		checkerror(err)

		err = gitlab.CreateRepo(repo)
		checkerror(err)

		err = git.CreateGitIgnoreFile()
		checkerror(err)

		err = git.PushRepo(gitlab.GitLabURL, gitlab.GitlabUser, basicUserInfo)
		if err != nil {
			removeFileErr := os.Remove(filepath.Join(configFolder, "git-push-gitlab"))
			if removeFileErr != nil {
				colorRed.Println("Error: " + removeFileErr.Error())
				os.Exit(0)
			}
		}
		checkerror(err)
	}
}
