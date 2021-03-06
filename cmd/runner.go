package main

import (
	"fmt"
	"log"
	"os"
	"talisman/checksumcalculator"
	"talisman/detector"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/prompt"
	"talisman/report"
	"talisman/scanner"
	"talisman/talismanrc"
	"talisman/utility"

	"github.com/spf13/afero"
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
	results   *helpers.DetectionResults
}

//NewRunner returns a new Runner.
func NewRunner(additions []gitrepo.Addition) *Runner {
	return &Runner{
		additions: additions,
		results:   helpers.NewDetectionResults(),
	}
}

//RunWithoutErrors will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
func (r *Runner) RunWithoutErrors(promptContext prompt.PromptContext) int {
	r.doRun()
	r.printReport(promptContext)
	return r.exitStatus()
}

//Scan scans git commit history for potential secrets and returns 0 or 1 as exit code
func (r *Runner) Scan(reportDirectory string, tRC *talismanrc.TalismanRC, ignoreHistory bool) int {
	fmt.Printf("\n\n")
	utility.CreateArt("Running Scan..")
	additions := scanner.GetAdditions(ignoreHistory)
	rcConfig := talismanrc.Get()
	setCustomSeverities(rcConfig)
	detector.DefaultChain(rcConfig).Test(additions, rcConfig, r.results)
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
	wd, _ := os.Getwd()
	repo := gitrepo.RepoLocatedAt(wd)
	gitTrackedFilesAsAdditions := repo.TrackedFilesAsAdditions()
	//Adding staged files for calculation
	gitTrackedFilesAsAdditions = append(gitTrackedFilesAsAdditions, repo.StagedAdditions()...)
	cc := checksumcalculator.NewChecksumCalculator(utility.DefaultSHA256Hasher{}, gitTrackedFilesAsAdditions)
	rcSuggestion := cc.SuggestTalismanRC(fileNamePatterns)
	if rcSuggestion != "" {
		fmt.Print(rcSuggestion)
		exitStatus = 0
	}
	return exitStatus
}

func (r *Runner) doRun() {
	rcConfig := talismanrc.Get()
	setCustomSeverities(rcConfig)
	scopeMap := getScopeConfig()
	additionsToScan := rcConfig.IgnoreAdditionsByScope(r.additions, scopeMap)
	detector.DefaultChain(rcConfig).Test(additionsToScan, rcConfig, r.results)
}

func setCustomSeverities(tRC *talismanrc.TalismanRC) {
	for _, cs := range tRC.CustomSeverities {
		severity.SeverityConfiguration[cs.Detector] = cs.Severity
	}
}

func getScopeConfig() map[string][]string {
	scopeConfig := map[string][]string{
		"node": {"yarn.lock", "package-lock.json", "node_modules/"},
		"go":   {"makefile", "go.mod", "go.sum", "Gopkg.toml", "Gopkg.lock", "glide.yaml", "glide.lock", "vendor/"},
	}
	return scopeConfig
}

func (r *Runner) printReport(promptContext prompt.PromptContext) {
	if r.results.HasWarnings() {
		fmt.Println(r.results.ReportWarnings())
	}
	if r.results.HasIgnores() || r.results.HasFailures() {
		fs := afero.NewOsFs()
		r.results.Report(fs, talismanrc.DefaultRCFileName, promptContext)
	}
}

func (r *Runner) exitStatus() int {
	if r.results.HasFailures() {
		return CompletedWithErrors
	}
	return CompletedSuccessfully
}
