package git

import (
	"github.com/urvil38/git-cli/lib"
	"github.com/urvil38/git-cli/questions"
	"gopkg.in/AlecAivazis/survey.v1"
)

func CreateGitIgnoreFile() error {
	gitignore := false
	err := survey.Ask(questions.WantsGitIgnore,&gitignore)
	if err != nil {
		return err
	}
	if gitignore {
		var gitignorefiles string
		err := survey.Ask(questions.CreateGitIgnore,&gitignorefiles)
		if err != nil {
			return err
		}
		f,err := lib.CreateFileInCurrentDir(".gitignore")
		defer f.Close()
		if err != nil {
			return err
		}
		_,err = f.WriteString(gitignorefiles)
		if err != nil {
			return err
		}
	}
	return nil
}