package main

import (
	"os"

	"talisman/gitrepo"
)

type PreCommitHook struct{
	runner
}

func NewPreCommitHook() *PreCommitHook {
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)

	return &PreCommitHook{*NewRunner(repo.GetDiffForStagedFiles())}
}
