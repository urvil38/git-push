package types

type BasicAuth struct {
	Username string `survey:"username"`
	Password string `survey:"password"`
}

type OAuth struct {
	Token string `survey:"token"`
}

type BasicUserInfo struct {
	Name string `survey:"name"`
	Email string `survey:"email"`
}

type Repo struct {
	RepoName string `survey:"reponame"`
	RepoDescription string `survey:"repodescription"`
	RepoType string `survey:"repotype"`
}

type RepoURL struct {
	HTMLURL string
	CloneURL string
	SSHURL string
}

type Service interface {
	Init() error
	CreateRepo(Repo) error
	CreateGitIgnoreFile() error
	PushRepo() error
}