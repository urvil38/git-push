package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gogit "gopkg.in/src-d/go-git.v4"
)

func GetCurrentWorkingDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	cwdarr := strings.Split(cwd, string(filepath.Separator))
	return cwdarr[len(cwdarr)-1:][0]
}

func CreateFileInCurrentDir(filename string) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := os.Create(cwd + string(filepath.Separator) + filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CheckIfFileIsExist(filename string) bool {
	cwd, err := os.Getwd()
	if _, err = os.Stat(cwd + string(filepath.Separator) + filename); err == nil {
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
