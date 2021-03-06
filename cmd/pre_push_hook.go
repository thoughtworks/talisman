package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"talisman/gitrepo"
)

const (
	//EmptySha represents the state of a brand new ref
	EmptySha string = "0000000000000000000000000000000000000000"
	//ShaId of the empty tree in Git
	EmptyTreeSha string = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
)

type PrePushHook struct {
	localRef, localCommit, remoteRef, remoteCommit string
}

func NewPrePushHook(localRef, localCommit, remoteRef, remoteCommit string) *PrePushHook {
	return &PrePushHook{localRef, localCommit, remoteRef, remoteCommit}
}

//If the outgoing ref does not exist on the remote, all commits on the local ref will be checked
//If the outgoing ref already exists, all additions in the range between "localSha" and "remoteSha" will be validated
func (p *PrePushHook) GetRepoAdditions() []gitrepo.Addition {
	if p.runningOnDeletedRef() {
		log.WithFields(log.Fields{
			"localRef":     p.localRef,
			"localCommit":  p.localCommit,
			"remoteRef":    p.remoteRef,
			"remoteCommit": p.remoteCommit,
		}).Info("Running on a deleted ref. Nothing to verify as outgoing changes are all deletions.")

		return []gitrepo.Addition{}
	}

	if p.runningOnNewRef() {
		log.WithFields(log.Fields{
			"localRef":     p.localRef,
			"localCommit":  p.localCommit,
			"remoteRef":    p.remoteRef,
			"remoteCommit": p.remoteCommit,
		}).Info("Running on a new ref. All changes in the ref will be verified.")

		return p.getRepoAdditionsFrom(EmptyTreeSha, p.localCommit)
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

func (p *PrePushHook) getRepoAdditions() []gitrepo.Addition {
	return p.getRepoAdditionsFrom(p.remoteCommit, p.localCommit)
}

func (p *PrePushHook) getRepoAdditionsFrom(oldCommit, newCommit string) []gitrepo.Addition {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	return repo.AdditionsWithinRange(oldCommit, newCommit)
}
