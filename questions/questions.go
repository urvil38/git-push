package questions

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

//ServiceNameQuestion defindes question for servicename
var ServiceNameQuestion = []*survey.Question{
	{
		Name: "servicename",
		Prompt: &survey.Select{
			Message: "Pleasse select a service you want to use",
			Options: []string{"Github","Bitbucket"},
			Default: "Github",
		},
	},
}