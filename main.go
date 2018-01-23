package main

import (
	"github.com/urvil38/git-cli/questions"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/urvil38/git-cli/color"
	"fmt"
)

var gitlogo = `
  ________.__  __                   .__  .__ 
 /  _____/|__|/  |_            ____ |  | |__|
/   \  ___|  \   __\  ______ _/ ___\|  | |  |
\    \_\  \  ||  |   /_____/ \  \___|  |_|  |
 \______  /__||__|            \___  >____/__|
        \/                        \/         
`

func main() {
	fmt.Println(color.Wrap(gitlogo))
	answer := struct {
		ServiceName string `survey:"servicename"`
	}{}

	err := survey.Ask(questions.ServiceNameQuestion,&answer)
	if err != nil {
		fmt.Println(err)
		return
	}
}
