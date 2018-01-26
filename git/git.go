package git

import (
	"fmt"
	"time"
	
	"github.com/urvil38/git-push/color"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/urvil38/git-push/utils"
	"gopkg.in/AlecAivazis/survey.v1"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func init() {
	remoteExists, err = utils.CheckRemoteRepo()
	if err != nil {
		return
	}
}

var (
	remoteExists bool
	err          error
)

func CreateGitIgnoreFile() error {
	gitignore := false
	if utils.CheckIfFileIsExist(".gitignore") {
		return nil
	}
	err := survey.Ask(questions.WantsGitIgnore, &gitignore)
	if err != nil {
		return err
	}
	if gitignore {
		var gitignorefiles string
		err := survey.Ask(questions.CreateGitIgnore, &gitignorefiles)
		if err != nil {
			return err
		}
		f, err := utils.CreateFileInCurrentDir(".gitignore")
		defer f.Close()
		if err != nil {
			return err
		}
		_, err = f.WriteString(gitignorefiles)
		if err != nil {
			return err
		}
	}
	return nil
}

func PushRepo(gitURL types.RepoURL, user types.User, basicUserInfo types.BasicUserInfo) error {
	if !remoteExists {
		r, err := gogit.PlainOpen(utils.GetCurrentWorkingDirPath())
		if err != nil {
			r,err = gogit.PlainInit(utils.GetCurrentWorkingDirPath(),false)
			if err != nil {
				return err
			}
		}
		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: gogit.DefaultRemoteName,
			URLs: []string{gitURL.CloneURL},
			
		})
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}
		err = utils.GitAddAll()
		if err != nil {
			return err
		}
		_, err = w.Commit("Initial Commit", &gogit.CommitOptions{
			Author: &object.Signature{
				Name:  basicUserInfo.Name,
				Email: basicUserInfo.Email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}

		auth := http.BasicAuth{
			Username: user.Username,
			Password: user.Password,
		}
		err = r.Push(&gogit.PushOptions{
			Auth: &auth,
		})
		if err != nil {
			return err
		}
		fmt.Println(color.Wrap("=> Successfully Pushed Repository", "FgGreen", "CrossedOut"))
		return nil
	}
	return nil
}
