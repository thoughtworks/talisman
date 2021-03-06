package main

import (
	"fmt"
	"github.com/spf13/afero"
	"talisman/detector"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/prompt"
	"talisman/talismanrc"
)

const (
	//CompletedSuccessfully is an exit status that says that the current runners run completed without errors
	CompletedSuccessfully int = 0

	//CompletedWithErrors is an exit status that says that the current runners run completed with failures
	CompletedWithErrors int = 1
)

//runner represents a single run of the validations for a given commit range
type runner struct {
	additions []gitrepo.Addition
	results   *helpers.DetectionResults
}

//NewRunner returns a new runner.
func NewRunner(additions []gitrepo.Addition) *runner {
	return &runner{
		additions: additions,
		results:   helpers.NewDetectionResults(),
	}
}

//Run will validate the commit range for errors and return either COMPLETED_SUCCESSFULLY or COMPLETED_WITH_ERRORS
func (r *runner) Run(tRC *talismanrc.TalismanRC, promptContext prompt.PromptContext) int {
	setCustomSeverities(tRC)
	additionsToScan := tRC.FilterAdditions(r.additions)
	detector.DefaultChain(tRC).Test(additionsToScan, tRC, r.results)
	r.printReport(promptContext)
	exitStatus := r.exitStatus()
	return exitStatus
}

func setCustomSeverities(tRC *talismanrc.TalismanRC) {
	for _, cs := range tRC.CustomSeverities {
		severity.SeverityConfiguration[cs.Detector] = cs.Severity
	}
}

func (r *runner) printReport(promptContext prompt.PromptContext) {
	if r.results.HasWarnings() {
		fmt.Println(r.results.ReportWarnings())
	}
	if r.results.HasIgnores() || r.results.HasFailures() {
		fs := afero.NewOsFs()
		r.results.Report(fs, talismanrc.DefaultRCFileName, promptContext)
	}
}

func (r *runner) exitStatus() int {
	if r.results.HasFailures() {
		return CompletedWithErrors
	}
	return CompletedSuccessfully
}
