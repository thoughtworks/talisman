package main

import (
	"fmt"
	"os"
	"talisman/checksumcalculator"
	"talisman/detector"
	"talisman/git_repo"
	"talisman/report"
	"talisman/scanner"
)

const (
	//CompletedSuccessfully is an exit status that says that the current runners run completed without errors
	CompletedSuccessfully int = 0

	//CompletedWithErrors is an exit status that says that the current runners run completed with failures
	CompletedWithErrors int = 1
)

//Runner represents a single run of the validations for a given commit range
type Runner struct {
	additions []git_repo.Addition
	results   *detector.DetectionResults
}

//NewRunner returns a new Runner.
func NewRunner(additions []git_repo.Addition) *Runner {
	return &Runner{additions, detector.NewDetectionResults()}
}

//RunWithoutErrors will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
func (r *Runner) RunWithoutErrors() int {
	r.doRun()
	r.printReport()
	return r.exitStatus()
}

//Scan scans git commit history for potential secrets and returns 0 or 1 as exit code
func (r *Runner) Scan() int {
	fmt.Println("Please wait while talisman scans entire repository including the git history...")
	additions := scanner.GetAdditions()
	ignores := detector.TalismanRCIgnore{}
	detector.DefaultChain().Test(additions, ignores, r.results)
	report.GenerateReport(r.results)
	fmt.Println("Please check report.html in your current directory for the talisman scan report")
	return r.exitStatus()
}

//RunChecksumCalculator runs the checksum calculator against the patterns given as input
func (r *Runner) RunChecksumCalculator(fileNamePatterns []string) int {
	exitStatus := 1
	cc := checksumcalculator.NewChecksumCalculator(fileNamePatterns)
	rcSuggestion := cc.SuggestTalismanRC()
	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		exitStatus = 0
	}
	return exitStatus
}

func (r *Runner) doRun() {
	ignoresNew := detector.ReadConfigFromRCFile(readRepoFile())
	detector.DefaultChain().Test(r.additions, ignoresNew, r.results)
}

func (r *Runner) printReport() {
	if r.results.HasWarnings() {
		fmt.Println(r.results.ReportWarnings())
	}
	if r.results.HasIgnores() || r.results.HasFailures() {
		fmt.Println(r.results.Report())
	}
}

func (r *Runner) exitStatus() int {
	if r.results.HasFailures() {
		return CompletedWithErrors
	}
	return CompletedSuccessfully
}

func readRepoFile() func(string) ([]byte, error) {
	wd, _ := os.Getwd()
	repo := git_repo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}
