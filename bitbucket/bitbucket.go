package bitbucket

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/briandowns/spinner"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	home = os.Getenv("HOME")
	configFilePath = filepath.Join(home, ".config", "git-push", "git-push-bitbucket")
	checkCredential()
}

func checkCredential() {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}
	b = encoding.Decode(string(b))
	credentials := strings.Split(string(b), "\n")
	BitbucketUser.Username = credentials[0]
	BitbucketUser.Password = credentials[1]
}

var (
	BitbucketUser  types.BasicAuth
	BitbuckerURL   types.RepoURL
	client         *bitbucket.Client
	configFilePath string
	home           string
	c              *color.Color
)

var (
	errNotGetHTMLURL  = errors.New("Not get right type of html url value from responsse")
	errNotGetCloneURL = errors.New("Not get right type of clone urls value from responsse")
)

func Init() error {

	if BitbucketUser.Username != "" || BitbucketUser.Password != "" {
		c.Println("=> You authenticated successfully ✓")
		return nil
	}
	err := survey.Ask(questions.BitbucketCredential, &BitbucketUser)
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
	client = bitbucket.NewBasicAuth(BitbucketUser.Username, BitbucketUser.Password)
	user, err := client.Users.Get(BitbucketUser.Username)
	if user == nil || err != nil {
		s.Stop()
		return errors.New("Invalid username or password ✗")
	}
	s.Stop()
	c.Println("=> You authenticated successfully ✓")

	b := new(bytes.Buffer)
	b.WriteString(BitbucketUser.Username + "\n" + BitbucketUser.Password)

	sEnc := encoding.Encode(b.Bytes())

	err = ioutil.WriteFile(configFilePath, []byte(sEnc), 0555)
	if err != nil {
		return err
	}
	return nil
}

func CreateRepo(repo types.Repo) error {
	client = bitbucket.NewBasicAuth(BitbucketUser.Username, BitbucketUser.Password)

	s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Fetching Repo URL from BitBucket ⚡"
	s.Start()

	r, err := client.Repositories.Repository.Create(&bitbucket.RepositoryOptions{
		Owner:       BitbucketUser.Username,
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
	c.Println("=> " + BitbuckerURL.HTMLURL)
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
	BitbuckerURL.HTMLURL = url
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
			BitbuckerURL.CloneURL, ok = v["href"].(string)
			checkErrCloneURL(ok)
		}
		if v["name"] == "ssh" {
			BitbuckerURL.SSHURL, ok = v["href"].(string)
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
