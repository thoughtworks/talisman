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
	ignoreEvaluator helpers.IgnoreEvaluator
	tRC             *talismanrc.TalismanRC
}

// Run scans git commit history for potential secrets and returns 0 or 1 as exit code
func (s *ScannerCmd) Run() int {
	fmt.Printf("\n\n")
	utility.CreateArt("Running Scan..")

	additionsToScan := s.tRC.FilterAdditions(s.additions)

	detector.DefaultChain(s.tRC, s.ignoreEvaluator).Test(additionsToScan, s.tRC, s.results)
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

// NewScannerCmd Returns a new scanner command
func NewScannerCmd(ignoreHistory bool, tRC *talismanrc.TalismanRC, reportDirectory string) *ScannerCmd {
	repoRoot, _ := os.Getwd()
	reader := gitrepo.NewBatchGitObjectHashReader(repoRoot)
	additions := scanner.GetAdditions(ignoreHistory, reader)
	ignoreEvaluator := helpers.ScanHistoryEvaluator()
	if ignoreHistory {
		ignoreEvaluator = helpers.BuildIgnoreEvaluator("default", tRC, gitrepo.RepoLocatedAt(repoRoot))
	}
	return &ScannerCmd{
		additions:       additions,
		results:         helpers.NewDetectionResults(talismanrc.ScanMode),
		reportDirectory: reportDirectory,
		ignoreEvaluator: ignoreEvaluator,
		tRC:             tRC,
	}
}
