package github

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/git"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/urvil38/git-push/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	configFilePath = utils.GetConfigFilePath("git-push-github")
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(b)
	credentials := strings.Split(string(b), "\n")
	GithubService.githubUser.Username = credentials[0]
	GithubService.githubUser.Password = credentials[1]
}

var (
	client         *github.Client
	configFilePath string
	c              *color.Color
	GithubService  githubservice
)

type githubservice struct {
	gitURL        types.RepoURL
	githubUser    types.BasicAuth
	basicUserInfo types.BasicUserInfo
}

//Init function ask for github username and password for basic auth
func (g githubservice) Init() error {
	if GithubService.githubUser.Username != "" || GithubService.githubUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}
	err := survey.Ask(questions.GithubCredential, &GithubService.githubUser)
	if err != nil {
		fmt.Println(err)
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Authenticating You ⚡"
	s.Start()

	err = authenticateUser(s)
	if err != nil {
		return err
	}
	return nil
}

func authenticateUser(s *spinner.Spinner) error {
	tp := github.BasicAuthTransport{
		Username: GithubService.githubUser.Username,
		Password: GithubService.githubUser.Password,
	}
	client = github.NewClient(tp.Client())
	ctx := context.Background()
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		s.Stop()
		return errors.New("Invalid username or password ✗")
	}
	s.Stop()
	c.Println("=> You authenticated successfully ✓")

	b := new(bytes.Buffer)
	b.WriteString(GithubService.githubUser.Username + "\n" + GithubService.githubUser.Password)

	sEnc := encoding.Encode(b.Bytes())

	err = ioutil.WriteFile(configFilePath, sEnc, 0555)
	if err != nil {
		return err
	}

	return nil
}

func (g githubservice) CreateRepo(repo types.Repo) error {
	tp := github.BasicAuthTransport{
		Username: GithubService.githubUser.Username,
		Password: GithubService.githubUser.Password,
	}
	client = github.NewClient(tp.Client())
	ctx := context.Background()
	r := &github.Repository{
		Name:        github.String(repo.RepoName),
		Description: github.String(repo.RepoDescription),
		Private:     github.Bool(repo.RepoType == "Private"),
	}

	s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Fetching Repo URL from Github ⚡"
	s.Start()

	repository, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		s.Stop()
		if strings.Contains(err.Error(), "exists") {
			return errors.New("Error: Same name of repository is exists on your account")
		}
		if strings.Contains(err.Error(), "private") {
			return errors.New("Error: Please upgrade your plan to create a new private repository")
		}
		return errors.New("Error while creating repository.Please check your internet connection ℹ")
	}

	stringify := func(str *string) string {
		return strings.Trim(github.Stringify(str), "\"")
	}

	GithubService.gitURL = types.RepoURL{
		HTMLURL:  stringify(repository.HTMLURL),
		CloneURL: stringify(repository.CloneURL),
		SSHURL:   stringify(repository.SSHURL),
	}

	s.Stop()
	c.Println("=> " + GithubService.gitURL.HTMLURL)
	return nil
}

func (g githubservice) CreateGitIgnoreFile() error {
	return git.CreateGitIgnoreFile()
}

func (g githubservice) PushRepo() error {
	var userConfigFile = utils.GetUserConfigFilePath()
	b, err := ioutil.ReadFile(userConfigFile)
	if err != nil {
		return err
	}
	userInfo := strings.Split(string(b), "\n")
	GithubService.basicUserInfo.Name = userInfo[0]
	GithubService.basicUserInfo.Email = userInfo[1]
	return git.PushRepo(GithubService.gitURL, GithubService.githubUser, GithubService.basicUserInfo)
}
