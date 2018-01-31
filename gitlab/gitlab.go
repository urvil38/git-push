package gitlab

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	home = os.Getenv("HOME")
	configFilePath = home + separator + ".config" + separator + "git-push" + separator + "git-push-gitlab"
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(string(b))
	credentials := strings.Split(string(b), "\n")
	GitlabToken.Token = credentials[0]
	GitlabUser.Username = credentials[1]
	GitlabUser.Password = credentials[2]
}

var (
	GitlabUser     types.BasicAuth
	GitlabToken    types.OAuth
	GitLabURL      types.RepoURL
	client         *gitlab.Client
	home           string
	configFilePath string
	c              *color.Color
)

const (
	separator = string(filepath.Separator)
)

func Init() error {

	if GitlabToken.Token != "" || GitlabUser.Username != "" || GitlabUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}

	if GitlabToken.Token == "" {
		err := survey.Ask(questions.GitlabToken, &GitlabToken)
		if err != nil {
			fmt.Println(err)
		}
	}

	user, err := checkToken()
	if err != nil {
		return err
	}

	if GitlabUser.Username == "" || GitlabUser.Password == "" {
		err := survey.Ask(questions.GitlabCredential, &GitlabUser)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = authenticateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func checkToken() (*gitlab.User, error) {
	client = gitlab.NewClient(nil, GitlabToken.Token)

	user, response, err := client.Users.CurrentUser()

	if err != nil {
		if response == nil {
			return nil, errors.New("ERROR: Please check your internet connection ℹ .")
		}
	}

	if err != nil || user == nil {
		return nil, errors.New("Invalid token ✗")
	}

	c.Println("=> Valid Token ✓")
	return user, nil
}

func authenticateUser(user *gitlab.User) error {

	if user.Username != GitlabUser.Username {
		return errors.New("Invalid username of password ✗")
	}

	c.Println("=> You authenticated successfully ✓")

	b := new(bytes.Buffer)
	b.WriteString(GitlabToken.Token + "\n" + GitlabUser.Username + "\n" + GitlabUser.Password)

	Estr := encoding.Encode(b.Bytes())

	err := ioutil.WriteFile(configFilePath, []byte(Estr), 0555)
	if err != nil {
		return err
	}

	return nil
}

func CreateRepo(repo types.Repo) error {
	client = gitlab.NewClient(nil, GitlabToken.Token)
	project, response, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.String(repo.RepoName),
		Description: gitlab.String(repo.RepoDescription),
		Visibility:  gitlab.Visibility(gitlab.VisibilityValue(strings.ToLower(repo.RepoType))),
	})
	if err != nil {
		if response != nil {
			if response.StatusCode == 400 {
				return errors.New("Error: Same name of repository is exists on your account")
			}
		}
		return errors.New("ERROR: Unable to create repo on gitlab.Please check your internet connection ℹ .")
	}

	GitLabURL.HTMLURL = project.WebURL
	GitLabURL.CloneURL = project.HTTPURLToRepo
	GitLabURL.SSHURL = project.SSHURLToRepo

	c.Println("=> " + GitLabURL.HTMLURL)

	return nil
}
