package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var gitPath = "tmp"
var loger = log.New(os.Stdout, "test_", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c, err := GetConfig()
	if err != nil {
		panic(err)
	}

	if !Exists(gitPath) {
		loger.Print("creat a temp dir to save git repos")
		err = os.Mkdir(gitPath, 0777)
		panic(err)
	}

	//var jsonConfig, err1 = json.Marshal(c)
	//if err1 == nil {
	//	fmt.Print(string(jsonConfig))
	//}

	for _, s := range c.Repos {
		sync(s, c.Config)
	}
}

func sync(config *syncConfig, globalConfig *globalConfig) {
	originSplit := strings.Split(config.Origin, "/")
	path := "./tmp/" + strings.Replace(originSplit[len(originSplit)-1], ".git", "", 1)

	fmt.Print(path)

	if !Exists(path) {
		loger.Printf("%s does not exist, clone from %s", path, config.Origin)

		cloneOptions := &git.CloneOptions{
			URL:          config.Origin,
			Progress:     os.Stdout,
			RemoteName:   config.OriginBranch,
			SingleBranch: true,
		}
		if config.OriginAuth.Group != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: config.OriginAuth.Username,
				Password: config.OriginAuth.Password,
			}
		}
		_, err := git.PlainClone(path, false, cloneOptions)
		if err != nil {
			loger.Fatal(err)
		}
	}

	worktree := &git.Worktree{}

	pullOptions := &git.PullOptions{
		RemoteName:   config.OriginBranch,
		SingleBranch: true,
	}
	if config.OriginAuth.Group != "" {
		pullOptions.Auth = &http.BasicAuth{
			Username: config.OriginAuth.Username,
			Password: config.OriginAuth.Password,
		}
	}

	worktree.Pull(pullOptions)

}
