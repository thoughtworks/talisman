package detector

import (
	"fmt"
	"os"
	"talisman/git_repo"
	"talisman/utility"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

type FailureData struct {
	FailuresInCommits map[string][]string
}

//DetectionResults represents all interesting information collected during a detection run.
//It serves as a collecting parameter for the tests performed by the various Detectors in the DetectorChain
//Currently, it keeps track of failures and ignored files.
//The results are grouped by FilePath for easy reporting of all detected problems with individual files.
type DetectionResults struct {
	failures map[git_repo.FilePath]*FailureData
	ignores  map[git_repo.FilePath][]string
	warnings map[git_repo.FilePath][]string
}

//NewDetectionResults is a new DetectionResults struct. It represents the pre-run state of a Detection run.
func NewDetectionResults() *DetectionResults {
	result := DetectionResults{make(map[git_repo.FilePath]*FailureData), make(map[git_repo.FilePath][]string), make(map[git_repo.FilePath][]string)}
	return &result
}

//Fail is used to mark the supplied FilePath as failing a detection for a supplied reason.
//Detectors are encouraged to provide context sensitive messages so that fixing the errors is made simple for the end user
//Fail may be called multiple times for each FilePath and the calls accumulate the provided reasons
func (r *DetectionResults) Fail(filePath git_repo.FilePath, message string, commits []string) {
	if r.failures[filePath] == nil {
		r.failures[filePath] = &FailureData{make(map[string][]string)}
	}
	if r.failures[filePath].FailuresInCommits == nil {
		r.failures[filePath].FailuresInCommits = make(map[string][]string)
	}
	existingCommits := r.failures[filePath].FailuresInCommits[message]
	if len(existingCommits) == 0 {
		r.failures[filePath].FailuresInCommits[message] = commits
	} else {
		r.failures[filePath].FailuresInCommits[message] = append(r.failures[filePath].FailuresInCommits[message], commits...)
	}
}

func (r *DetectionResults) Warn(filePath git_repo.FilePath, message string) {
	warnings, ok := r.warnings[filePath]
	if !ok {
		r.warnings[filePath] = []string{message}
	} else {
		r.warnings[filePath] = append(warnings, message)
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

//HasFailures answers if any Failures were detected for any FilePath in the current run
func (r *DetectionResults) HasFailures() bool {
	return len(r.failures) > 0
}

//HasIgnores answers if any FilePaths were ignored in the current run
func (r *DetectionResults) HasIgnores() bool {
	return len(r.ignores) > 0
}

func (r *DetectionResults) HasWarnings() bool {
	return len(r.warnings) > 0
}

func (r *DetectionResults) HasDetectionMessages() bool {
	return r.HasWarnings() || r.HasFailures() || r.HasIgnores()
}

//Successful answers if no detector was able to find any possible result to fail the run
func (r *DetectionResults) Successful() bool {
	return !r.HasFailures()
}

//GetFailures returns the various reasons that a given FilePath was marked as failing by all the detectors in the current run
func (r *DetectionResults) GetFailures(fileName git_repo.FilePath) *FailureData {
	return r.failures[fileName]
}

func (r *DetectionResults) ReportWarnings() string {
	var result string
	var filePathsForWarnings []string
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Warnings"})
	table.SetRowLine(true)

	for filePath := range r.warnings {
		filePathsForWarnings = append(filePathsForWarnings, string(filePath))
		warningData := r.ReportFileWarnings(filePath)
		data = append(data, warningData...)
	}
	filePathsForWarnings = unique(filePathsForWarnings)
	if len(r.warnings) > 0 {
		table.AppendBulk(data)
		table.Render()
		result = result + fmt.Sprintf("\n\x1b[33mPlease review the above file(s) to make sure that no sensitive content is being pushed\x1b[0m\n")
		result = result + fmt.Sprintf("\n")
	}
	return result
}

//Report returns a string documenting the various failures and ignored files for the current run
func (r *DetectionResults) Report() string {
	var result string
	var filePathsForIgnoresAndFailures []string
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Errors"})
	table.SetRowLine(true)

	for filePath := range r.failures {
		filePathsForIgnoresAndFailures = append(filePathsForIgnoresAndFailures, string(filePath))
		failureData := r.ReportFileFailures(filePath)
		data = append(data, failureData...)
	}

	for filePath := range r.ignores {
		filePathsForIgnoresAndFailures = append(filePathsForIgnoresAndFailures, string(filePath))
		// ignoreData := r.ReportFileIgnores(filePath)
		// data = append(data, ignoreData...)
	}
	filePathsForIgnoresAndFailures = unique(filePathsForIgnoresAndFailures)
	if len(r.failures) > 0 {
		fmt.Printf("\n\x1b[1m\x1b[31mTalisman Report:\x1b[0m\x1b[0m\n")
		table.AppendBulk(data)
		table.Render()
		result = result + fmt.Sprintf("\n\x1b[33mIf you are absolutely sure that you want to ignore the above files from talisman detectors, consider pasting the following format in .talismanrc file in the project root\x1b[0m\n")
		result = result + r.suggestTalismanRC(filePathsForIgnoresAndFailures)
		result = result + fmt.Sprintf("\n\n")
	}
	return result
}

func (r *DetectionResults) suggestTalismanRC(filePaths []string) string {
	var fileIgnoreConfigs []FileIgnoreConfig
	for _, filePath := range filePaths {
		currentChecksum := utility.CollectiveSHA256Hash([]string{filePath})
		fileIgnoreConfig := FileIgnoreConfig{filePath, currentChecksum, []string{}}
		fileIgnoreConfigs = append(fileIgnoreConfigs, fileIgnoreConfig)
	}

	talismanRcIgnoreConfig := TalismanRCIgnore{fileIgnoreConfigs}
	m, _ := yaml.Marshal(&talismanRcIgnoreConfig)
	return string(m)
}

//ReportFileFailures adds a string to table documenting the various failures detected on the supplied FilePath by all detectors in the current run
func (r *DetectionResults) ReportFileFailures(filePath git_repo.FilePath) [][]string {
	failures := r.failures[filePath]
	var data [][]string
	if len(failures.FailuresInCommits) > 0 {
		for failureMessage := range failures.FailuresInCommits {
			if len(failureMessage) > 150 {
				failureMessage = failureMessage[:150] + "\n" + failureMessage[150:]
			}
			data = append(data, []string{string(filePath), failureMessage})
		}
	}
	return data
}

func (r *DetectionResults) ReportFileWarnings(filePath git_repo.FilePath) [][]string {
	warnings := r.warnings[filePath]
	var data [][]string
	if len(warnings) > 0 {
		for _, warning := range warnings {
			data = append(data, []string{string(filePath), warning})
		}
	}
	return data
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

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NewFailureData() FailureData {
	return FailureData{make(map[string][]string)}
}
