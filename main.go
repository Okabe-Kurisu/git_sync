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
	if CheckIfError(err) {
		return
	}

	if !Exists(gitPath) {
		loger.Print("creat a temp dir to save git repos")
		err = os.Mkdir(gitPath, 0777)
		if CheckIfError(err) {
			return
		}
	}

	cron := crontab.New()
	for _, s := range c.Repos {
		SyncConfig := &syncConfig{}
		err := DeepCopy(SyncConfig, s)
		if CheckIfError(err) {
			return
		}
		sync(SyncConfig)
		err = cron.AddFunc(s.Frequency, func() {
			sync(SyncConfig)
		})
		if CheckIfError(err) {
			return
		}
	}
	cron.Run()
}

func sync(config *syncConfig) {
	originUrl, targetUrl := strings.Replace(config.Origin, ".git", "/tree/", 1), strings.Replace(config.Target, ".git", "/tree/", 1)
	log.Printf("%s%s is syncing to %s%s", originUrl, config.OriginBranch, config.Target, targetUrl)
	remoteName := config.OriginBranch + config.TargetBranch
	repo := getRepository(config, remoteName+"pull")

	log.Printf("pulling from %s%s", originUrl, config.OriginBranch)
	pullRemote(config, repo, remoteName+"pull")

	log.Printf("pushing to %s%s", targetUrl, config.TargetBranch)
	pushRepository(config, repo, remoteName+"push")

	log.Printf("%s%s has already synced to %s%s", originUrl, config.OriginBranch, targetUrl, config.TargetBranch)
}

func pullRemote(config *syncConfig, repo *git.Repository, remoteName string) {
	worktree, err := repo.Worktree()
	if CheckIfError(err) {
		return
	}

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
		if CheckIfError(err) {
			return
		}
	}
}

func pushRepository(config *syncConfig, repo *git.Repository, remoteName string) {

	remote, err := repo.CreateRemote(&gitConfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{config.Target},
	})
	if err != git.ErrRemoteExists {
		if CheckIfError(err) {
			return
		}
	}
	remote, err = repo.Remote(remoteName)
	if CheckIfError(err) {
		return
	}

	pushOptions := &git.PushOptions{
		RemoteName: remoteName,
		Progress:   os.Stdout,
		RefSpecs: []gitConfig.RefSpec{
			gitConfig.RefSpec("+refs/heads/" + config.OriginBranch + ":refs/heads/" + config.TargetBranch),
		},
	}
	if config.TargetAuth.Group != "" {
		pushOptions.Auth = &http.BasicAuth{
			Username: config.TargetAuth.Username,
			Password: config.TargetAuth.Password,
		}
	}
	if config.IsSyncTags {
		pushOptions.RefSpecs = append(pushOptions.RefSpecs, "refs/tags/*:refs/tags/*")
	}

	err = remote.Push(pushOptions)
	if err != git.NoErrAlreadyUpToDate {
		if CheckIfError(err) {
			return
		}
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
		repos, err = git.PlainClone(path, false, cloneOptions)

	} else {
		repos, err = git.PlainOpen(path)
	}
	if CheckIfError(err) {
		return nil
	}
	return repos
}
