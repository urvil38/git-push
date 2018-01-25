package github

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urvil38/git-cli/color"
	"github.com/urvil38/git-cli/encoding"

	"github.com/google/go-github/github"
	"github.com/urvil38/git-cli/questions"
	"github.com/urvil38/git-cli/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	Home = os.Getenv("HOME")
}

var (
	Home   string
	client *github.Client
	GitURL types.RepoURL
)

//Init function ask for github username and password for basic auth
func Init(answer *types.Answer) error {
	if answer.User.Username != "" || answer.User.Password != "" {
		fmt.Println(color.Wrap("=> You authenticated successfully", "FgGreen", "CrossedOut"))
		return nil
	}
	err := survey.Ask(questions.GithubCredential, &answer.User)
	if err != nil {
		fmt.Println(err)
	}

	err = authenticateUser(answer)
	if err != nil {
		return err
	}
	return nil
}

func authenticateUser(answer *types.Answer) error {
	tp := github.BasicAuthTransport{
		Username: answer.User.Username,
		Password: answer.User.Password,
	}
	client = github.NewClient(tp.Client())
	ctx := context.Background()
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return errors.New("Invalid username or password")
	}

	fmt.Println(color.Wrap("=> You authenticated successfully", "FgGreen", "CrossedOut"))

	f, err := os.Create(Home + "/.config/git-cli/git-cli")
	defer f.Close()

	b := new(bytes.Buffer)
	b.WriteString(answer.Username + "\n")
	b.WriteString(answer.Password)

	sEnc := encoding.Encode(b.Bytes())

	_, err = f.Write([]byte(sEnc))
	if err != nil {
		return err
	}

	return nil
}

func CreateRepo(answer *types.Answer) error {
	tp := github.BasicAuthTransport{
		Username: answer.Username,
		Password: answer.Password,
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
		return errors.New("Error while creating repository")
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
