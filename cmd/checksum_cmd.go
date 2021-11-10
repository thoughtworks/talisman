package main

import (
	"fmt"
	"os"
	"talisman/checksumcalculator"
	"talisman/gitrepo"
	"talisman/utility"

	"github.com/sirupsen/logrus"
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
	hasher := utility.MakeHasher("checksum", wd)

	err := hasher.Start()
	if err != nil {
		logrus.Errorf("unable to start hasher: %v", err)
		return EXIT_FAILURE
	}

	cc := checksumcalculator.NewChecksumCalculator(hasher, gitTrackedFilesAsAdditions)
	hasher.Shutdown()

	rcSuggestion := cc.SuggestTalismanRC(s.fileNamePatterns)

	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		return EXIT_SUCCESS
	}
	return EXIT_FAILURE
}
