package main

import (
	"fmt"
	"os"
	"talisman/detector"
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/report"
	"talisman/scanner"
	"talisman/talismanrc"
	"talisman/utility"

	logr "github.com/sirupsen/logrus"
)

const (
	SCAN_MODE = "scan"
)

type ScannerCmd struct {
	additions       []gitrepo.Addition
	results         *helpers.DetectionResults
	reportDirectory string
}

//Run scans git commit history for potential secrets and returns 0 or 1 as exit code
func (s *ScannerCmd) Run(tRC *talismanrc.TalismanRC) int {
	fmt.Printf("\n\n")
	utility.CreateArt("Running ScanMode..")
	detector.DefaultChain(tRC, "default").Test(s.additions, tRC, s.results)
	reportsPath, err := report.GenerateReport(s.results, s.reportDirectory)
	if err != nil {
		logr.Errorf("error while generating report: %v", err)
		return EXIT_FAILURE
	}

	fmt.Printf("\nPlease check '%s' folder for the talisman scan report\n\n", reportsPath)
	return s.exitStatus()
}

func (s *ScannerCmd) exitStatus() int {
	if s.results.HasFailures() {
		return EXIT_FAILURE
	}
	return EXIT_SUCCESS
}

//NewScannerCmd Returns a new scanner command
func NewScannerCmd(ignoreHistory bool, reportDirectory string) *ScannerCmd {
	repoRoot, _ := os.Getwd()
	reader := gitrepo.NewBatchGitObjectHashReader(repoRoot)
	additions := scanner.GetAdditions(ignoreHistory, reader)
	return &ScannerCmd{
		additions:       additions,
		results:         helpers.NewDetectionResults(talismanrc.ScanMode),
		reportDirectory: reportDirectory,
	}
}
