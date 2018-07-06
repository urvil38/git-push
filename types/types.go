//Package types represents types used by git-push
package types

//BasicAuth contains Username and password for authentication with client
type BasicAuth struct {
	Username string `survey:"username"`
	Password string `survey:"password"`
}

//OAuth token for authenticating with gitlab.
//Because gitlab required token and also basicAuth like username and password.
type OAuth struct {
	Token string `survey:"token"`
}

//BasicUserInfo contains information for commit messages with git
type BasicUserInfo struct {
	Name string `survey:"name"`
	Email string `survey:"email"`
}

//Repo represent git repository information
//RepoType can be "public" or "private"
type Repo struct {
	RepoName string `survey:"reponame"`
	RepoDescription string `survey:"repodescription"`
	RepoType string `survey:"repotype"`
}

//RepoURL contains URL for remote git repository
//It can be html , clone and ssh urls
//Ex: html:  https://github.com/urvil38/git-push ,for HTML page which user can directly open in browser
//Ex: clone: https://github.com/urvil38/git-push.git ,setup remote repo as cloneURL created by 
//Ex: ssh:   git@github.com:urvil38/git-push.git ,right now ssh is not supported as authentication mechanism
//TODO: Support ssh for authentication mechanism
type RepoURL struct {
	HTMLURL string
	CloneURL string
	SSHURL string
}

//Service interface for all things to do in order to push repository in selected provider
type Service interface {
	//Authenticate with client using methods they support.
	//Github and Bitbucket require username and password as authentication mechanism
	//Gitlab requires username and password as well as OAuth token in order to authenticate 
	Authenticate() error

	//CreateRepo creates empty repository using appropriate client choosen by user
	//Ex:If user select github then It uses github client to create empty repository based on information given by CreateRepo() function
	CreateRepo(Repo) error

	//It ask user if he/she want's to add .gitignore file
	//If current folder has already .gitignore file this function do not ask for he/she want to ignore
	CreateGitIgnoreFile() error

	//PushRepo push current folder to remote repository on selected service
	//Ex: git init -> git add . -> git commit -> git push
	PushRepo() error
}