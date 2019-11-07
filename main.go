package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
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
		err := sync(s, c.Config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sync(config *syncConfig, globalConfig *globalConfig) (err error) {
	log.Printf("%s:%s is syncing to %s:%s", config.Origin, config.OriginBranch, config.Target, config.TargetBranch)
	repo, err := getRepository(config)
	CheckIfError(err)

	worktree, err := repo.Worktree()
	CheckIfError(err)

	pullOptions := &git.PullOptions{
		RemoteName:   config.OriginBranch,
		SingleBranch: true,
		Force:        true,
	}
	if config.OriginAuth.Group != "" {
		pullOptions.Auth = &http.BasicAuth{
			Username: config.OriginAuth.Username,
			Password: config.OriginAuth.Password,
		}
	}

	log.Printf("pulling from %s:%s", config.Origin, config.OriginBranch)
	err = worktree.Pull(pullOptions)
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	log.Printf("pushing to %s:%s", config.Target, config.TargetBranch)

	remoteName := config.OriginBranch + config.TargetBranch
	remote, err := repo.CreateRemote(&gitConfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{config.Target},
	})
	if err != git.ErrRemoteExists {
		CheckIfError(err)
	}
	remote, err = repo.Remote(remoteName)
	CheckIfError(err)

	pushOptions := &git.PushOptions{
		RemoteName: remoteName,
		Prune:      config.IsForce,
		RefSpecs: []gitConfig.RefSpec{
			gitConfig.RefSpec("refs/heads/" + config.OriginBranch + ":refs/heads/" + config.TargetBranch),
		},
	}
	if config.TargetAuth.Group != "" {
		pushOptions.Auth = &http.BasicAuth{
			Username: config.TargetAuth.Username,
			Password: config.TargetAuth.Password,
		}
	}

	err = remote.Push(pushOptions)
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	log.Printf("%s:%s has already synced to %s:%s", config.Origin, config.OriginBranch, config.Target, config.TargetBranch)
	return nil
}

func getRepository(config *syncConfig) (*git.Repository, error) {
	originSplit := strings.Split(config.Origin, "/")
	path := "./tmp/" + strings.Replace(originSplit[len(originSplit)-1], ".git", "", 1)
	var repos *git.Repository
	var err error
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
		repos, err = git.PlainClone(path, false, cloneOptions)

	} else {
		repos, err = git.PlainOpen(path)
	}
	if err != nil {
		return nil, err
	} else {
		return repos, nil
	}
}
