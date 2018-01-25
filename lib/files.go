package lib

import (
	"path/filepath"
	"os"
	"strings"
)

func GetCurrentWorkingDirName() string {
	cwd,err := os.Getwd()
	if err != nil {
		return ""
	}
	cwdarr := strings.Split(cwd,string(filepath.Separator))
	return cwdarr[len(cwdarr)-1:][0]
}

func CreateFileInCurrentDir(filename string) (*os.File,error) {
	cwd,err := os.Getwd()
	if err != nil {
		return nil,err
	}
	f ,err := os.Create(cwd + string(filepath.Separator) + filename)
	if err != nil {
		return nil,err
	}
	return f,nil
}