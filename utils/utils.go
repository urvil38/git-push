package utils

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"path/filepath"

	gogit "gopkg.in/src-d/go-git.v4"
)

const help = `
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

func GetConfigFolderPath() string {
	home := getHomeEnv("HOME")
	return filepath.Join(home, ".config", "git-push")
}

func GetConfigFilePath(fileName string) string {
	home := getHomeEnv("HOME")
	return filepath.Join(home, ".config", "git-push", fileName)
}

func GetUserConfigFilePath() string {
	home := getHomeEnv("HOME")
	return filepath.Join(home, ".config", "git-push", "userInfo")
}

func getHomeEnv(env string) string {
	home := os.Getenv("HOME")
	if home == "" {
		fmt.Println(help)
		os.Exit(0)
	}
	return home
}

func CreateDir(dirName string) error {
	return os.MkdirAll(dirName,0777)
}

func GetCurrentWorkingDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	cwdarr := strings.Split(cwd,string(filepath.Separator))
	return cwdarr[len(cwdarr)-1:][0]
}

func CreateFileInCurrentDir(filename string) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := os.Create(filepath.Join(cwd,filename))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CheckIfFileIsExist(filename string) bool {
	cwd, err := os.Getwd()
	if _, err = os.Stat(filepath.Join(cwd,filename)); err == nil {
		return true
	}
	return false
}

func GetCurrentWorkingDirPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "couldn't get current dir path"
	}
	return cwd
}

func CheckRemoteRepo() (bool, error) {
	r, err := gogit.PlainOpen(GetCurrentWorkingDirPath())
	if err != nil {
		return false, err
	}
	remotes, err := r.Remotes()
	if err != nil {
		return false, err
	}
	if len(remotes) < 1 {
		return false, nil
	}
	return true,nil
}

func GitAddAll() error {
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ResetAccount(reset string) error {
	if reset == "all" {
		err := os.RemoveAll(GetConfigFolderPath())
		if err != nil {
			return err
		}
		return nil
	}

	filename := "git-push-"+reset
	err := os.Remove(GetConfigFilePath(filename))
	if err != nil {
		return err
	}
	return nil
}
