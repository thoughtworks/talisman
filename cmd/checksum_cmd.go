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
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)
	cc := checksumcalculator.NewChecksumCalculator(utility.MakeHasher("checksum", wd), gitTrackedFilesAsAdditions)
	rcSuggestion := cc.SuggestTalismanRC(s.fileNamePatterns)
	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		return EXIT_SUCCESS
	}
	return EXIT_FAILURE
}
