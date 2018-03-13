package bitbucket

import (
	"github.com/urvil38/git-push/git"
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urvil38/git-push/utils"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/briandowns/spinner"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	configFilePath = utils.GetConfigFilePath()
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(string(b))
	credentials := strings.Split(string(b), "\n")
	BitbucketService.bitbucketUser.Username = credentials[0]
	BitbucketService.bitbucketUser.Password = credentials[1]
}

var (
	client         *bitbucket.Client
	configFilePath string
	home           string
	c              *color.Color
	BitbucketService bitbucketService
)

var (
	errNotGetHTMLURL  = errors.New("Not get right type of html url value from responsse")
	errNotGetCloneURL = errors.New("Not get right type of clone urls value from responsse")
)

type bitbucketService struct {
	bitbucketURL types.RepoURL
	bitbucketUser types.BasicAuth
	basicUserInfo types.BasicUserInfo
}

func (b bitbucketService) Init() error {

	if BitbucketService.bitbucketUser.Username != "" || BitbucketService.bitbucketUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}
	err := survey.Ask(questions.BitbucketCredential, &BitbucketService.bitbucketUser)
	if err != nil {
		return err
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
	client = bitbucket.NewBasicAuth(BitbucketService.bitbucketUser.Username, BitbucketService.bitbucketUser.Password)
	user, err := client.Users.Get(BitbucketService.bitbucketUser.Username)
	if user == nil || err != nil {
		s.Stop()
		return errors.New("Invalid username or password ✗")
	}
	s.Stop()
	c.Println("=> You authenticated successfully ✓")

	bytes := new(bytes.Buffer)
	bytes.WriteString(BitbucketService.bitbucketUser.Username + "\n" + BitbucketService.bitbucketUser.Password)

	sEnc := encoding.Encode(bytes.Bytes())

	err = ioutil.WriteFile(configFilePath, []byte(sEnc), 0555)
	if err != nil {
		return err
	}
	return nil
}

func (b bitbucketService) CreateRepo(repo types.Repo) error {
	client = bitbucket.NewBasicAuth(BitbucketService.bitbucketUser.Username, BitbucketService.bitbucketUser.Password)

	s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Fetching Repo URL from BitBucket ⚡"
	s.Start()

	r, err := client.Repositories.Repository.Create(&bitbucket.RepositoryOptions{
		Owner:       BitbucketService.bitbucketUser.Username,
		Repo_slug:   repo.RepoName,
		Description: repo.RepoDescription,
		Is_private:  formatBool(repo.RepoType == "Private"),
	})
	if err != nil {
		s.Stop()
		return errors.New("Error:Couldn't create repository.Make sure you don't have same repository name on bitbucket or Please check your internet connection ℹ")
	}
	err = typeCheckHTMLURL(r)
	if err != nil {
		s.Stop()
		return err
	}
	err = typeCheckCloneURL(r)
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()
	c.Println("=> " + BitbucketService.bitbucketURL.HTMLURL)
	return nil
}

func checkErrHTMLURL(ok bool) error {
	if !ok {
		return errNotGetHTMLURL
	}
	return nil
}

func checkErrCloneURL(ok bool) error {
	if !ok {
		return errNotGetCloneURL
	}
	return nil
}

func typeCheckHTMLURL(r *bitbucket.Repository) error {
	value := r.Links["html"]
	urls, ok := value.(map[string]interface{})
	checkErrHTMLURL(ok)
	url, ok := urls["href"].(string)
	checkErrHTMLURL(ok)
	BitbucketService.bitbucketURL.HTMLURL = url
	return nil
}

func typeCheckCloneURL(r *bitbucket.Repository) error {
	value := r.Links["clone"]
	values, ok := value.([]interface{})
	checkErrCloneURL(ok)
	for _, value := range values {
		v, ok := value.(map[string]interface{})
		checkErrCloneURL(ok)
		if v["name"] == "https" {
			BitbucketService.bitbucketURL.CloneURL, ok = v["href"].(string)
			checkErrCloneURL(ok)
		}
		if v["name"] == "ssh" {
			BitbucketService.bitbucketURL.SSHURL, ok = v["href"].(string)
			checkErrCloneURL(ok)
		}
	}
	return nil
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func (b bitbucketService) CreateGitIgnoreFile() error {
	return git.CreateGitIgnoreFile()
}

func (b bitbucketService) PushRepo() error {
	var userConfigFile = filepath.Join(home, ".config", "git-push", "userInfo")
	bytes, err := ioutil.ReadFile(userConfigFile)
	if err != nil {
		return err
	}
	userInfo := strings.Split(string(bytes), "\n")
	BitbucketService.basicUserInfo.Name = userInfo[0]
	BitbucketService.basicUserInfo.Email = userInfo[1]
	return git.PushRepo(BitbucketService.bitbucketURL,BitbucketService.bitbucketUser,BitbucketService.basicUserInfo)
}