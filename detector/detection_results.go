package detector

import (
	"fmt"
	"strings"

	"github.com/thoughtworks/talisman/git_repo"
)

//DetectionResults represents all interesting information collected during a detection run.
//It serves as a collecting parameter for the tests performed by the various Detectors in the DetectorChain
//Currently, it keeps track of failures and ignored files.
//The results are grouped by FilePath for easy reporting of all detected problems with individual files.
type DetectionResults struct {
	failures map[git_repo.FilePath][]string
	ignores  map[git_repo.FilePath][]string
}

//NewDetectionResults is a new DetectionResults struct. It represents the pre-run state of a Detection run.
func NewDetectionResults() *DetectionResults {
	result := DetectionResults{make(map[git_repo.FilePath][]string), make(map[git_repo.FilePath][]string)}
	return &result
}

//Fail is used to mark the supplied FilePath as failing a detection for a supplied reason.
//Detectors are encouraged to provide context sensitive messages so that fixing the errors is made simple for the end user
//Fail may be called multiple times for each FilePath and the calls accumulate the provided reasons
func (r *DetectionResults) Fail(filePath git_repo.FilePath, message string) {
	errors, ok := r.failures[filePath]
	if !ok {
		r.failures[filePath] = []string{message}
	} else {
		r.failures[filePath] = append(errors, message)
	}
}

//Ignore is used to mark the supplied FilePath as being ignored.
//The most common reason for this is that the FilePath is Denied by the Ignores supplied to the Detector, however, Detectors may use more sophisticated reasons to ignore files.
func (r *DetectionResults) Ignore(filePath git_repo.FilePath, detector string) {
	ignores, ok := r.ignores[filePath]
	if !ok {
		r.ignores[filePath] = []string{detector}
	} else {
		r.ignores[filePath] = append(ignores, detector)
	}
}

//HasFailures answers if any failures were detected for any FilePath in the current run
func (r *DetectionResults) HasFailures() bool {
	return len(r.failures) > 0
}

//HasIgnores answers if any FilePaths were ignored in the current run
func (r *DetectionResults) HasIgnores() bool {
	return len(r.ignores) > 0
}

//Successful answers if no detector was able to find any possible result to fail the run
func (r *DetectionResults) Successful() bool {
	return !r.HasFailures()
}

//Failures returns the various reasons that a given FilePath was marked as failing by all the detectors in the current run
func (r *DetectionResults) Failures(fileName git_repo.FilePath) []string {
	return r.failures[fileName]
}

//Report returns a string documenting the various failures and ignored files for the current run
func (r *DetectionResults) Report() string {
	var result string
	for filePath := range r.failures {
		result = result + r.ReportFileFailures(filePath)
	}

	if len(r.ignores) > 0 {
		result = result + fmt.Sprintf("The following files were ignored:\n")
	}
	for filePath := range r.ignores {
		result = result + fmt.Sprintf("\t%s was ignored by .talismanignore for the following detectors: %s\n", filePath, strings.Join(r.ignores[filePath], ", "))
	}
	return result
}

//ReportFileFailures returns a string documenting the various failures detected on the supplied FilePath by all detectors in the current run
func (r *DetectionResults) ReportFileFailures(filePath git_repo.FilePath) string {
	failures := r.failures[filePath]
	if len(failures) > 0 {
		result := fmt.Sprintf("The following errors were detected in %s\n", filePath)
		for _, failure := range failures {
			result = result + fmt.Sprintf("\t %s\n", failure)
		}
		return result
	}
	return ""
}

func (r *DetectionResults) failurePaths() []git_repo.FilePath {
	return keys(r.failures)
}

func (r *DetectionResults) ignorePaths() []git_repo.FilePath {
	return keys(r.ignores)
}

func keys(aMap map[git_repo.FilePath][]string) []git_repo.FilePath {
	var result []git_repo.FilePath
	for filePath := range aMap {
		result = append(result, filePath)
	}
	return result
}
