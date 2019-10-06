package main

import (
	"os"

	"talisman/gitrepo"
)

type PreCommitHook struct{}

func NewPreCommitHook() *PreCommitHook {
	return &PreCommitHook{}
}

func (p *PreCommitHook) GetRepoAdditions() []gitrepo.Addition {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	return repo.GetDiffForStagedFiles()
}
