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
	hasher           utility.SHA256Hasher
	repoRoot         string
}

func NewChecksumCmd(fileNamePatterns []string) *ChecksumCmd {
	wd, _ := os.Getwd()
	hasher := utility.MakeHasher("checksum", wd)
	return &ChecksumCmd{fileNamePatterns: fileNamePatterns, hasher: hasher, repoRoot: wd}
}

func (s *ChecksumCmd) Run() int {
	repo := gitrepo.RepoLocatedAt(s.repoRoot)
	err := s.hasher.Start()
	if err != nil {
		logrus.Errorf("unable to start hasher: %v", err)
		return EXIT_FAILURE
	}

	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)

	cc := checksumcalculator.NewChecksumCalculator(s.hasher, gitTrackedFilesAsAdditions)
	rcSuggestion := cc.SuggestTalismanRC(s.fileNamePatterns)
	s.hasher.Shutdown()

	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		return EXIT_SUCCESS
	}
	return EXIT_FAILURE
}
