package main

import (
	"github.com/urvil38/git-push/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urvil38/git-push/bitbucket"
	"github.com/urvil38/git-push/github"
	"github.com/urvil38/git-push/gitlab"
	"github.com/urvil38/git-push/types"
)

func service(service types.Service,repo types.Repo) error {
	if err := service.Init(); err != nil {
		return err
	}
	
	if err := service.CreateRepo(repo); err != nil {
		return err
	}

	if err := service.CreateGitIgnoreFile(); err != nil {
		return err
	}

	return service.PushRepo()
}

func invokeService(serviceName string,repo types.Repo) error {
	switch serviceName {
	case "GitHub":
		return service(github.GithubService,repo)
	case "BitBucket":
		return service(bitbucket.BitbucketService,repo)
	case "GitLab":
		err := service(gitlab.GitlabService,repo)
		if err.Error() == "giterror" {
			fmt.Printf("%s\n", red("Error: "+"Please check gitlab username or password are correct"))
			removeFileErr := os.Remove(filepath.Join(utils.GetConfigFolderPath(), "git-push-gitlab"))
			if removeFileErr != nil {
				fmt.Printf("%s\n", red("Error: "+removeFileErr.Error()))
				os.Exit(0)
			}
			os.Exit(1)
		}
		return err
	}
	return nil
}
