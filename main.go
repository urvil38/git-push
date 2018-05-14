package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/urvil38/git-push/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	red = color.New(color.FgRed, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	green = color.New(color.FgGreen, color.Bold).SprintFunc()

	reset := flag.String("reset", "", "Use for Resetting Account, Equivalent to Logout\nExample: git-push -reset [github | bitbucket | gitlab | all]")
	flag.Parse()
	if *reset != "" {
		reset := strings.ToLower(*reset)
		err := utils.ResetAccount(reset)
		if err != nil {
			_, ok := err.(*os.PathError)
			if ok {
				fmt.Printf("%s\n", red("Cound't Reset Account.You are not Logged in to "+reset+" Account"))
			}
			os.Exit(0)
		}
		fmt.Printf("%s\n", green("Successfully Reset "+reset+" Account"))
		os.Exit(0)
	}

	userConfigFile = utils.GetUserConfigFilePath()
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
	green          func(...interface{}) string
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

	err = invokeService(serviceName, repo)
	if err != nil {
		fmt.Printf("%s\n", red("=> "+err.Error()))
		os.Exit(0)
	}
}
