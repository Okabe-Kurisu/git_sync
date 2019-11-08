package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"

	crontab "github.com/robfig/cron"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var gitPath = "tmp"
var loger = log.New(os.Stdout, "test_", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c, err := GetConfig()
	CheckIfError(err)

	if !Exists(gitPath) {
		loger.Print("creat a temp dir to save git repos")
		err = os.Mkdir(gitPath, 0777)
		CheckIfError(err)
	}

	cron := crontab.New()
	for _, s := range c.Repos {
		err := cron.AddFunc(s.Frequency, func() {
			sync(s)
		})
		CheckIfError(err)
	}
	cron.Run()
}

func sync(config *syncConfig) {
	log.Printf("%s:%s is syncing to %s:%s", config.Origin, config.OriginBranch, config.Target, config.TargetBranch)
	remoteName := config.OriginBranch + config.TargetBranch
	repo := getRepository(config, remoteName+"pull")

	log.Printf("pulling from %s:%s", config.Origin, config.OriginBranch)
	pullRemote(config, repo, remoteName+"pull")

	log.Printf("pushing to %s:%s", config.Target, config.TargetBranch)
	pushRepository(config, repo, remoteName+"push")

	log.Printf("%s:%s has already synced to %s:%s", config.Origin, config.OriginBranch, config.Target, config.TargetBranch)
}

func pullRemote(config *syncConfig, repo *git.Repository, remoteName string) {
	worktree, err := repo.Worktree()
	CheckIfError(err)

	pullOptions := &git.PullOptions{
		RemoteName:    remoteName,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + config.OriginBranch),
		SingleBranch:  true,
		Force:         true,
		Progress:      os.Stdout,
	}
	if config.OriginAuth.Group != "" {
		pullOptions.Auth = &http.BasicAuth{
			Username: config.OriginAuth.Username,
			Password: config.OriginAuth.Password,
		}
	}

	err = worktree.Pull(pullOptions)
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}
}

func pushRepository(config *syncConfig, repo *git.Repository, remoteName string) {

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
		Progress:   os.Stdout,
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

}

func getRepository(config *syncConfig, remoteName string) *git.Repository {
	originSplit := strings.Split(config.Origin, "/")
	path := "./tmp/" + strings.Replace(originSplit[len(originSplit)-1], ".git", "", 1)
	var repos *git.Repository
	var err error
	if !Exists(path) {
		loger.Printf("%s does not exist, clone from %s", path, config.Origin)

		cloneOptions := &git.CloneOptions{
			URL:           config.Origin,
			Progress:      os.Stdout,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + config.OriginBranch),
			RemoteName:    remoteName,
			SingleBranch:  true,
		}
		if config.OriginAuth.Group != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: config.OriginAuth.Username,
				Password: config.OriginAuth.Password,
			}
		}
		log.Printf("%s", plumbing.ReferenceName("refs/heads/"+config.OriginBranch))
		repos, err = git.PlainClone(path, false, cloneOptions)

	} else {
		repos, err = git.PlainOpen(path)
	}
	CheckIfError(err)
	return repos
}
