package questions

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/urvil38/git-cli/lib"
)

//ServiceNameQuestion defindes question for servicename
var ServiceName = []*survey.Question{
	{
		Name: "servicename",
		Prompt: &survey.Select{
			Message: "Pleasse select a service you want to use:",
			Options: []string{"Github","Bitbucket"},
			Default: "Github",
		},
	},
}
//GithubCredential ask for username and password for basic auth
var GithubCredential = []*survey.Question{
	{
		Name: "username",
		Prompt: &survey.Input{
			Message: "Enter your username or Email:",
			Help: "Please give your github username or email address",
		},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Enter your Password:",
		},
		Validate: survey.Required,
	},
}

var GithubRepoInfo = []*survey.Question{
	{
		Name: "reponame",
		Prompt: &survey.Input{
			Message: "Enter name of the repository:",
			Default: lib.GetCurrentWorkingDirName(),
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
			Options: []string{"Public","Private"},
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
			Message: "Please add files you want to ignore for git",
			Default: "node_modules\n*.gem\n*.rbc",
			HideDefault: true,
			AppendDefault: true,
		},
	},
}