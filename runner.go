package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/badrij/talisman/git_repo"
)

const (
	//EMPTY_SHA represents the state of a brand new ref
	EMPTY_SHA string = "0000000000000000000000000000000000000000"

	//COMPLETED_SUCCESSFULLY is an exit status that says that the current runners run completed without errors
	COMPLETED_SUCCESSFULLY int = 0

	//COMPLETED_WITH_ERRORS is an exit status that says that the current runners run completed with failures
	COMPLETED_WITH_ERRORS int = 1
)

//Runner represents a single run of the validations for a given commit range
type Runner struct {
	localRef, localCommit, remoteRef, remoteCommit string
	results                                        *DetectionResults
}

//NewRunner returns a new Runner.
func NewRunner(localRef, localCommit, remoteRef, remoteCommit string) *Runner {
	return &Runner{localRef, localCommit, remoteRef, remoteCommit, NewDetectionResults()}
}

//RunWithoutErrors will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
//Brand new repositoris are not validated at all
//If the outgoing ref does not exist on the remote, all commits on the local ref will be checked
//If the outgoing ref already exists, all additions in the range beween "localSha" and "remoteSha" will be validated
func (r *Runner) RunWithoutErrors() int {
	if r.runningOnDeletedRef() {
		log.WithFields(log.Fields{
			"localRef":     r.localRef,
			"localCommit":  r.localCommit,
			"remoteRef":    r.remoteRef,
			"remoteCommit": r.remoteCommit,
		}).Info("Running on a deleted ref. Nothing to verify as outgoing changes are all deletions.")
		return COMPLETED_SUCCESSFULLY
	}
	if r.runningOnNewRef() {
		return r.checkAllCommitsInNewRef()
	}
	return r.checkAllCommitsInRange()
}

func (r *Runner) checkAllCommitsInNewRef() int {
	log.WithFields(log.Fields{
		"localRef":     r.localRef,
		"localCommit":  r.localCommit,
		"remoteRef":    r.remoteRef,
		"remoteCommit": r.remoteCommit,
	}).Info("Running on a new ref. All changes in the ref will be verified.")
	return COMPLETED_SUCCESSFULLY
}

func (r *Runner) checkAllCommitsInRange() int {
	log.WithFields(log.Fields{
		"localRef":     r.localRef,
		"localCommit":  r.localCommit,
		"remoteRef":    r.remoteRef,
		"remoteCommit": r.remoteCommit,
	}).Info("Running on an existing ref. All changes in the commit range will be verified.")
	r.doRun()
	r.printReport()
	return r.exitStatus()
}

func (r *Runner) doRun() {
	ignores := ReadIgnoresFromFile(readRepoFile())
	DefaultDetectorChain().Test(r.getRepoAdditions(), ignores, r.results)
}

func (r *Runner) printReport() {
	if r.results.HasIgnores() || r.results.HasFailures() {
		fmt.Println(r.results.Report())
	}
}

func (r *Runner) exitStatus() int {
	if r.results.HasFailures() {
		return COMPLETED_WITH_ERRORS
	}
	return COMPLETED_SUCCESSFULLY
}

func (r *Runner) getRepoAdditions() []git_repo.Addition {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.Additions(r.remoteCommit, r.localCommit)
}

func (r *Runner) runningOnDeletedRef() bool {
	return r.localCommit == EMPTY_SHA
}

func (r *Runner) runningOnNewRef() bool {
	return r.remoteCommit == EMPTY_SHA
}

func readRepoFile() func(string) ([]byte, error) {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}
