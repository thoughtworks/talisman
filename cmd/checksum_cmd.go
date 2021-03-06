package main

import (
	"fmt"
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/utility"
)

type ChecksumCmd struct {
	fileNamePatterns []string
}

func NewChecksumCmd(fileNamePatterns []string) *ChecksumCmd {
	return &ChecksumCmd{fileNamePatterns: fileNamePatterns}
}

func (s *ChecksumCmd) Run() int {
	exitStatus := 1
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)
	cc := checksumcalculator.NewChecksumCalculator(utility.DefaultSHA256Hasher{}, gitTrackedFilesAsAdditions)
	rcSuggestion := cc.SuggestTalismanRC(s.fileNamePatterns)
	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		exitStatus = 0
	}
	return exitStatus
}
