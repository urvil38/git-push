package bitbucket

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/urvil38/git-push/encoding"
	"github.com/urvil38/git-push/questions"
	"github.com/urvil38/git-push/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

func init() {
	c = color.New(color.FgGreen, color.Bold)
	home = os.Getenv("HOME")
	configFilePath = home + separator + ".config" + separator + "git-push" + separator + "git-push-bitbucket"
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
	BitbucketUser  types.User
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

const (
	separator = string(filepath.Separator)
)

func Init() error {

	if BitbucketUser.Username != "" || BitbucketUser.Password != "" {
		c.Println("=> You authenticated successfully")
		return nil
	}
	err := survey.Ask(questions.BitbucketCredential, &BitbucketUser)
	if err != nil {
		return err
	}
	err = authenticateUser()
	if err != nil {
		return err
	}
	return nil
}

func authenticateUser() error {
	client = bitbucket.NewBasicAuth(BitbucketUser.Username, BitbucketUser.Password)
	user, err := client.Users.Get(BitbucketUser.Username)
	if user == nil || err != nil {
		return errors.New("Invalid username or password")
	}
	b := new(bytes.Buffer)
	b.WriteString(BitbucketUser.Username + "\n")
	b.WriteString(BitbucketUser.Password)

	sEnc := encoding.Encode(b.Bytes())

	err = ioutil.WriteFile(configFilePath, []byte(sEnc), 0555)
	if err != nil {
		return err
	}
	return nil
}

func CreateRepo(answer types.Answer) error {
	client = bitbucket.NewBasicAuth(BitbucketUser.Username, BitbucketUser.Password)
	r, err := client.Repositories.Repository.Create(&bitbucket.RepositoryOptions{
		Owner:       BitbucketUser.Username,
		Repo_slug:   answer.RepoName,
		Description: answer.RepoDescription,
		Is_private:  formatBool(answer.RepoType == "Private"),
	})
	if err != nil {
		return errors.New("Error:Couldn't create repository.Please check your internet connection or make sure you don't have same repository name on bitbucket")
	}
	typeCheckHTMLURL(r)
	typeCheckCloneURL(r)
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

func typeCheckCloneURL(r *bitbucket.Repository) {
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
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
