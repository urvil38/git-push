package questions

import (
	"github.com/urvil38/git-push/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

//ServiceName defindes question for servicename
var ServiceName = []*survey.Question{
	{
		Name: "servicename",
		Prompt: &survey.Select{
			Message: "Pleasse select a service you want to use:",
			Options: []string{"GitHub", "BitBucket", "GitLab"},
			Default: "GitHub",
		},
	},
}

//GithubCredential ask for username and password for basic auth
var GithubCredential = []*survey.Question{
	{
		Name: "username",
		Prompt: &survey.Input{
			Message: "Enter your github username or Email:",
			Help:    "Please give your github username or email address",
		},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Enter your github Password:",
		},
		Validate: survey.Required,
	},
}

var BitbucketCredential = []*survey.Question{
	{
		Name: "username",
		Prompt: &survey.Input{
			Message: "Enter your bitbucket Username:",
			Help:    "Please give your bitbucket Username",
		},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Enter your bitbucket Password:",
		},
		Validate: survey.Required,
	},
}

var GitlabCredential = []*survey.Question{
	{
		Name: "username",
		Prompt: &survey.Input{
			Message: "Enter your GitLab Username:",
			Help:    "Please give your GitLab Username",
		},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Enter your GitLab Password:",
		},
		Validate: survey.Required,
	},
}

var GitlabToken = []*survey.Question{
	{
		Name: "token",
		Prompt: &survey.Password{
			Message: "Enter your GitLab Token:",
			Help:    "You need to provide gitlab OAuth token here",
		},
		Validate: survey.Required,
	},
}

var UserInfo = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Enter Your name:",
			Help:    "Please give your name for configure git commit",
		},
		Validate: survey.Required,
	},
	{
		Name: "email",
		Prompt: &survey.Input{
			Message: "Enter your Email:",
			Help:    "Please give your email for configure git commit",
		},
		Validate: survey.Required,
	},
}

var GithubRepoInfo = []*survey.Question{
	{
		Name: "reponame",
		Prompt: &survey.Input{
			Message: "Enter name of the repository:",
			Default: utils.GetCurrentWorkingDirName(),
		},
		Validate: survey.Required,
	},
	{
		Name: "repodescription",
		Prompt: &survey.Input{
			Message: "Enter description of the repository:",
		},
	},
	{
		Name: "repotype",
		Prompt: &survey.Select{
			Message: "Public or Private:",
			Options: []string{"Public", "Private"},
			Default: "Public",
		},
	},
}

var WantsGitIgnore = []*survey.Question{
	{
		Name: "gitignore",
		Prompt: &survey.Confirm{
			Message: "You want to add .gitignore file?",
		},
	},
}

var CreateGitIgnore = []*survey.Question{
	{
		Name: "gitignorefile",
		Prompt: &survey.Editor{
			Message:       "Please add files you want to ignore for git",
			Default:       "node_modules\n*.gem\n*.rbc\n.vscode\n.idea\n",
			HideDefault:   true,
			AppendDefault: true,
		},
	},
}