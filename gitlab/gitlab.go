package gitlab

import (
	"github.com/urvil38/git-push/utils"
	"github.com/urvil38/git-push/git"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/briandowns/spinner"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	configFilePath = utils.GetConfigFilePath("git-push-gitlab")
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(b)
	credentials := strings.Split(string(b), "\n")
	GitlabService.gitlabToken.Token = credentials[0]
	GitlabService.gitlabUser.Username = credentials[1]
	GitlabService.gitlabUser.Password = credentials[2]
}

var (
	client         *gitlab.Client
	configFilePath string
	c              *color.Color
	GitlabService gitlabService
)

type gitlabService struct {
	gitlabUser types.BasicAuth
	gitlabToken types.OAuth
	basicUserInfo types.BasicUserInfo
	gitlabURL types.RepoURL
}

func (g gitlabService) Init() error {

	if GitlabService.gitlabToken.Token != "" || GitlabService.gitlabUser.Username != "" || GitlabService.gitlabUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}

	if GitlabService.gitlabToken.Token == "" {
		err := survey.Ask(questions.GitlabToken, &GitlabService.gitlabToken)
		if err != nil {
			fmt.Println(err)
		}
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Authenticating You ⚡"
	s.Start()

	user, err := checkToken(s)
	if err != nil {
		return err
	}

	if GitlabService.gitlabUser.Username == "" || GitlabService.gitlabUser.Password == "" {
		err := survey.Ask(questions.GitlabCredential, &GitlabService.gitlabUser)
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

func checkToken(s *spinner.Spinner) (*gitlab.User, error) {
	client = gitlab.NewClient(nil, GitlabService.gitlabToken.Token)

	user, response, err := client.Users.CurrentUser()

	if err != nil {
		s.Stop()
		if response == nil {
			return nil, errors.New("ERROR: Please check your internet connection ℹ ")
		}
	}

	if err != nil || user == nil {
		s.Stop()
		return nil, errors.New("Invalid token ✗")
	}
	s.Stop()
	c.Println("=> Valid Token ✓")
	return user, nil
}

func authenticateUser(user *gitlab.User) error {

	if user.Username != GitlabService.gitlabUser.Username {
		return errors.New("Invalid username of password ✗")
	}

	c.Println("=> You authenticated successfully ✓")

	b := new(bytes.Buffer)
	b.WriteString(GitlabService.gitlabToken.Token + "\n" + GitlabService.gitlabUser.Username + "\n" + GitlabService.gitlabUser.Password)

	Estr := encoding.Encode(b.Bytes())

	err := ioutil.WriteFile(configFilePath, Estr, 0555)
	if err != nil {
		return err
	}

	return nil
}

func (g gitlabService) CreateRepo(repo types.Repo) error {
	client = gitlab.NewClient(nil, GitlabService.gitlabToken.Token)

	s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Fetching Repo URL from GitLab ⚡"
	s.Start()

	project, response, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.String(repo.RepoName),
		Description: gitlab.String(repo.RepoDescription),
		Visibility:  gitlab.Visibility(gitlab.VisibilityValue(strings.ToLower(repo.RepoType))),
	})
	if err != nil {
		s.Stop()
		if response != nil {
			if response.StatusCode == 400 {
				return errors.New("Error: Same name of repository is exists on your account")
			}
		}
		return errors.New("ERROR: Unable to create repo on gitlab.Please check your internet connection ℹ .")
	}

	GitlabService.gitlabURL.HTMLURL = project.WebURL
	GitlabService.gitlabURL.CloneURL = project.HTTPURLToRepo
	GitlabService.gitlabURL.SSHURL = project.SSHURLToRepo

	s.Stop()
	c.Println("=> " + GitlabService.gitlabURL.HTMLURL)

	return nil
}

func (g gitlabService) CreateGitIgnoreFile() error {
	return git.CreateGitIgnoreFile()
}

func (g gitlabService) PushRepo() error {
	var userConfigFile = utils.GetUserConfigFilePath()
	b, err := ioutil.ReadFile(userConfigFile)
	if err != nil {
		return err
	}
	userInfo := strings.Split(string(b), "\n")
	GitlabService.basicUserInfo.Name = userInfo[0]
	GitlabService.basicUserInfo.Email = userInfo[1]
	err = git.PushRepo(GitlabService.gitlabURL,GitlabService.gitlabUser,GitlabService.basicUserInfo)
	if err != nil {
		return errors.New("giterror")
	}
	return nil
} 