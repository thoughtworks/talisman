package main

import (
	"fmt"
	"talisman/detector"
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/report"
	"talisman/scanner"
	"talisman/talismanrc"
	"talisman/utility"

	logr "github.com/sirupsen/logrus"
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
	detector.DefaultChain(tRC, PrePush).Test(s.additions, tRC, s.results)
	reportsPath, err := report.GenerateReport(s.results, s.reportDirectory)
	if err != nil {
		logr.Errorf("error while generating report: %v", err)
		return CompletedWithErrors
	}

	fmt.Printf("\nPlease check '%s' folder for the talisman scan report\n\n", reportsPath)
	return s.exitStatus()
}

func (s *ScannerCmd) exitStatus() int {
	if s.results.HasFailures() {
		return CompletedWithErrors
	}
	return CompletedSuccessfully
}

//NewScannerCmd Returns a new scanner command
func NewScannerCmd(ignoreHistory bool, reportDirectory string) *ScannerCmd {
	additions := scanner.GetAdditions(ignoreHistory)
	return &ScannerCmd{
		additions:       additions,
		results:         helpers.NewDetectionResults(talismanrc.ScanMode),
		reportDirectory: reportDirectory,
	}
}
