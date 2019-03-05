package main

import (
	"os"

	"talisman/git_repo"
)

type PreCommitHook struct{}

func NewPreCommitHook() *PreCommitHook {
	return &PreCommitHook{}
}

func (p *PreCommitHook) GetRepoAdditions() []git_repo.Addition {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.GetDiffForStagedFiles()
}
