package git

import (
	"errors"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
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
	c = color.New(color.FgGreen, color.Bold)
	remoteExists, err = utils.CheckRemoteRepo()
	if err != nil {
		return
	}
}

var (
	remoteExists bool
	err          error
	c            *color.Color
)

func CreateGitIgnoreFile() error {
	var gitignore bool
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

func PushRepo(gitURL types.RepoURL, user types.BasicAuth, basicUserInfo types.BasicUserInfo) error {
	if !remoteExists {
		r, err := gogit.PlainOpen(utils.GetCurrentWorkingDirPath())
		if err != nil {
			r, err = gogit.PlainInit(utils.GetCurrentWorkingDirPath(), false)
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

		s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
		s.Color("yellow")
		s.Suffix = " Working hard to pushing your Repository ⚡"
		s.Start()

		err = r.Push(&gogit.PushOptions{
			Auth: &auth,
		})
		if err != nil {
			s.Stop()
			return errors.New("ERROR: Unable to push repository.Please check your username or password are correct ℹ")
		}

		s.Stop()
		c.Println("=> Successfully Pushed Repository ✓")
		return nil
	}
	return nil
}
