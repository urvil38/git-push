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

	"github.com/google/go-github/github"
	"github.com/urvil38/git-push/color"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	home = os.Getenv("HOME")
	configFilePath = home + separator + ".config" + separator + "git-push" + separator + "git-push-github"
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
	GithubUser     types.User
)

const (
	separator = string(filepath.Separator)
)

//Init function ask for github username and password for basic auth
func Init() error {
	if GithubUser.Username != "" || GithubUser.Password != "" {
		fmt.Println(color.Wrap("=> You authenticated successfully", "FgGreen", "CrossedOut"))
		return nil
	}
	err := survey.Ask(questions.GithubCredential, &GithubUser)
	if err != nil {
		fmt.Println(err)
	}

	err = authenticateUser()
	if err != nil {
		return err
	}
	return nil
}

func authenticateUser() error {
	tp := github.BasicAuthTransport{
		Username: GithubUser.Username,
		Password: GithubUser.Password,
	}
	client = github.NewClient(tp.Client())
	ctx := context.Background()
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return errors.New("Invalid username or password")
	}

	fmt.Println(color.Wrap("=> You authenticated successfully", "FgGreen", "CrossedOut"))
	
	b := new(bytes.Buffer)
	b.WriteString(GithubUser.Username + "\n")
	b.WriteString(GithubUser.Password)
	
	sEnc := encoding.Encode(b.Bytes())
	
	err = ioutil.WriteFile(configFilePath,[]byte(sEnc),0555)
	if err != nil {
		return err
	}

	return nil
}

func CreateRepo(answer types.Answer) error {
	tp := github.BasicAuthTransport{
		Username: GithubUser.Username,
		Password: GithubUser.Password,
	}
	client := github.NewClient(tp.Client())
	ctx := context.Background()
	repo := &github.Repository{
		Name:        github.String(answer.RepoName),
		Description: github.String(answer.RepoDescription),
		Private:     github.Bool(answer.RepoType == "Private"),
	}
	repository, response, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		if response != nil && response.StatusCode == 422 {
			return errors.New("Error: " + "Same name of repository is exists on your account")
		}
		return errors.New("Error while creating repository.Please check your internet connection")
	}
	stringify := func(str *string) string {
		return strings.Trim(github.Stringify(str), "\"")
	}
	GitURL = types.RepoURL{
		HTMLURL:  stringify(repository.HTMLURL),
		CloneURL: stringify(repository.CloneURL),
		SSHURL:   stringify(repository.SSHURL),
	}
	fmt.Println(color.Wrap("=> "+GitURL.HTMLURL, "FgGreen", "CrossedOut"))
	return nil
}
