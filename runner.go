package main

import (
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
	"talisman/checksumcalculator"
	"talisman/detector"
	"talisman/gitrepo"
	"talisman/prompt"
	"talisman/report"
	"talisman/scanner"
	"talisman/utility"
)

const (
	//CompletedSuccessfully is an exit status that says that the current runners run completed without errors
	CompletedSuccessfully int = 0

	//CompletedWithErrors is an exit status that says that the current runners run completed with failures
	CompletedWithErrors int = 1
)

//Runner represents a single run of the validations for a given commit range
type Runner struct {
	additions []gitrepo.Addition
	results   *detector.DetectionResults
}

//NewRunner returns a new Runner.
func NewRunner(additions []gitrepo.Addition) *Runner {
	return &Runner{
		additions: additions,
		results:   detector.NewDetectionResults(),
	}
}

//RunWithoutErrors will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
func (r *Runner) RunWithoutErrors(prompter prompt.Prompt) int {
	r.doRun()
	r.printReport(prompter)
	return r.exitStatus()
}

//Scan scans git commit history for potential secrets and returns 0 or 1 as exit code
func (r *Runner) Scan(reportDirectory string) int {

	fmt.Printf("\n\n")
	utility.CreateArt("Running Scan..")
	additions := scanner.GetAdditions()
	ignores := detector.TalismanRCIgnore{}
	detector.DefaultChain().Test(additions, ignores, r.results)
	reportsPath, err := report.GenerateReport(r.results, reportDirectory)
	if err != nil {
		log.Printf("error while generating report: %v", err)
		return CompletedWithErrors
	}
	fmt.Printf("\nPlease check '%s' folder for the talisman scan report\n", reportsPath)
	fmt.Printf("\n")
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
	rcConfigIgnores := detector.ReadConfigFromRCFile(readRepoFile())
	scopeMap := getScopeConfig()
	additionsToScan := detector.IgnoreAdditionsByScope(r.additions, rcConfigIgnores, scopeMap);
	detector.DefaultChain().Test(additionsToScan, rcConfigIgnores, r.results)
}

func getScopeConfig() map[string][]string {
	scopeConfig := map[string][]string{
		"node": {"yarn.lock", "package-lock.json", "node_modules/"},
		"go":   {"makefile", "go.mod", "go.sum", "Gopkg.toml", "Gopkg.lock", "glide.yaml", "glide.lock", "vendor/"},
	}
	return scopeConfig
}

func (r *Runner) printReport(prompter prompt.Prompt) {
	if r.results.HasWarnings() {
		fmt.Println(r.results.ReportWarnings())
	}
	if r.results.HasIgnores() || r.results.HasFailures() {
		fs := afero.NewOsFs()
		r.results.Report(fs, detector.DefaultRCFileName, prompter)
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
	repo := gitrepo.RepoLocatedAt(wd)
	return repo.ReadRepoFileOrNothing
}
