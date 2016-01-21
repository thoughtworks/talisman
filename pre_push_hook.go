package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
)

const (
	//EmptySha represents the state of a brand new ref
	EmptySha string = "0000000000000000000000000000000000000000"
)

type PrePushHook struct {
	localRef, localCommit, remoteRef, remoteCommit string
}

func NewPrePushHook(localRef, localCommit, remoteRef, remoteCommit string) *PrePushHook {
	return &PrePushHook{localRef, localCommit, remoteRef, remoteCommit}
}

//Brand new repositoris are not validated at all
//If the outgoing ref does not exist on the remote, all commits on the local ref will be checked
//If the outgoing ref already exists, all additions in the range beween "localSha" and "remoteSha" will be validated
func (p *PrePushHook) GetRepoAdditions() []git_repo.Addition {
	if p.runningOnDeletedRef() {
		log.WithFields(log.Fields{
			"localRef":     p.localRef,
			"localCommit":  p.localCommit,
			"remoteRef":    p.remoteRef,
			"remoteCommit": p.remoteCommit,
		}).Info("Running on a deleted ref. Nothing to verify as outgoing changes are all deletions.")

		return []git_repo.Addition{}
	}

	if p.runningOnNewRef() {
		log.WithFields(log.Fields{
			"localRef":     p.localRef,
			"localCommit":  p.localCommit,
			"remoteRef":    p.remoteRef,
			"remoteCommit": p.remoteCommit,
		}).Info("Running on a new ref. All changes in the ref will be verified.")

		return []git_repo.Addition{}
	}

	log.WithFields(log.Fields{
		"localRef":     p.localRef,
		"localCommit":  p.localCommit,
		"remoteRef":    p.remoteRef,
		"remoteCommit": p.remoteCommit,
	}).Info("Running on an existing ref. All changes in the commit range will be verified.")

	return p.getRepoAdditions()
}

func (p *PrePushHook) runningOnDeletedRef() bool {
	return p.localCommit == EmptySha
}

func (p *PrePushHook) runningOnNewRef() bool {
	return p.remoteCommit == EmptySha
}

func (p *PrePushHook) getRepoAdditions() []git_repo.Addition {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.AdditionsWithinRange(p.remoteCommit, p.localCommit)
}
