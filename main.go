package main

import (
	"github.com/git-cli/logo"
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
	fmt.Print(logo.Wrap(gitlogo))
}
