package types

type Answer struct {
	ServiceName string `survey:"servicename"`
	Repo
}

type User struct {
	Username string `survey:"username"`
	Password string `survey:"password"`
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