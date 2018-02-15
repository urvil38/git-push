package github

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/briandowns/spinner"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	home = os.Getenv("HOME")
	configFilePath = filepath.Join(home, ".config", "git-push", "git-push-github")
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(string(b))
	credentials := strings.Split(string(b), "\n")
	GithubUser.Username = credentials[0]
	GithubUser.Password = credentials[1]
}

var (
	home           string
	client         *github.Client
	configFilePath string
	GitURL         types.RepoURL
	GithubUser     types.BasicAuth
	c              *color.Color
)

//Init function ask for github username and password for basic auth
func Init() error {
	if GithubUser.Username != "" || GithubUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}
	err := survey.Ask(questions.GithubCredential, &GithubUser)
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
		Username: GithubUser.Username,
		Password: GithubUser.Password,
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
	b.WriteString(GithubUser.Username + "\n" + GithubUser.Password)

	sEnc := encoding.Encode(b.Bytes())

	err = ioutil.WriteFile(configFilePath, []byte(sEnc), 0555)
	if err != nil {
		return err
	}

	return nil
}

func CreateRepo(repo types.Repo) error {
	tp := github.BasicAuthTransport{
		Username: GithubUser.Username,
		Password: GithubUser.Password,
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
		return errors.New("Error while creating repository.Please check your internet connection❗")
	}
	stringify := func(str *string) string {
		return strings.Trim(github.Stringify(str), "\"")
	}
	GitURL = types.RepoURL{
		HTMLURL:  stringify(repository.HTMLURL),
		CloneURL: stringify(repository.CloneURL),
		SSHURL:   stringify(repository.SSHURL),
	}
	s.Stop()
	c.Println("=> " + GitURL.HTMLURL)
	return nil
}
