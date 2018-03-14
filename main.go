package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/utils"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	red = color.New(color.FgRed, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()

	userConfigFile = utils.GetConfigFilePath()
	configFolder = utils.GetConfigFolderPath()

	err := utils.CreateDir(configFolder)
	if err != nil {
		log.Fatal(err)
	}

	checkUserInfo()

	remoteExists, _ = utils.CheckRemoteRepo()
	if remoteExists {
		fmt.Printf("%s%s%s\n",
			red("Sorry, this tool will not help you because working repository is already on github or bitbucket or gitlab!\nℹ You can use "),
			yellow("$ git push origin master"),
			red(" to push changes."))
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

var (
	basicUserInfo  types.BasicUserInfo
	remoteExists   bool
	err            error
	home           string
	userConfigFile string
	configFolder   string
	version        string
	red            func(...interface{}) string
	yellow         func(...interface{}) string
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
)

func main() {

	fmt.Printf("%s\n", yellow(fmt.Sprintf(banner, version)))

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

	err := invokeService(serviceName,repo)
	if err != nil {
		fmt.Printf("%s\n", red("=> "+err.Error()))
		os.Exit(0)
	}
}
