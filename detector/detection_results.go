package detector

import (
	"fmt"
	"os"
	"strconv"
	"talisman/git_repo"

	"github.com/olekukonko/tablewriter"
	yaml "gopkg.in/yaml.v2"
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

func (r *DetectionResults) ScannerReport() string {
	var result string
	var filePathsForFailures []string
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Errors", "Number of commits"})
	table.SetRowLine(true)

	for filePath := range r.failures {
		filePathsForFailures = append(filePathsForFailures, string(filePath))
		failureData := r.ReportFileFailuresForScanner(filePath)
		data = append(data, failureData...)
	}

	filePathsForFailures = unique(filePathsForFailures)
	if len(r.failures) > 0 {
		fmt.Printf("\n\x1b[1m\x1b[31mTalisman Report:\x1b[0m\x1b[0m\n")
		table.AppendBulk(data)
		table.Render()
	}
	return result
}

func (r *DetectionResults) suggestTalismanRC(filePaths []string) string {
	var fileIgnoreConfigs []FileIgnoreConfig
	for _, filePath := range filePaths {
		currentChecksum := CalculateCollectiveHash([]string{filePath})
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
	if len(failures) > 0 {
		for _, failure := range failures {
			if len(failure) > 150 {
				failure = failure[:150] + "\n" + failure[150:]
			}
			data = append(data, []string{string(filePath), failure})
		}
	}
	return data
}

func (r *DetectionResults) ReportFileFailuresForScanner(filePath git_repo.FilePath) [][]string {
	failures := r.failures[filePath]
	var data [][]string
	var failuresWithoutDuplicates []string
	if len(failures) > 0 {
		duplicate_frequency := make(map[string]int)
		for _, failure := range failures {
			_, exist := duplicate_frequency[failure]
			if exist {
				duplicate_frequency[failure] += 1
			} else {
				if len(failure) > 150 {
					failure = failure[:150] + "\n" + failure[150:]
				}
				duplicate_frequency[failure] = 1
				failuresWithoutDuplicates = append(failuresWithoutDuplicates, failure)
			}
		}
		for _, failure := range failuresWithoutDuplicates {
			numberOfCommits := strconv.Itoa(duplicate_frequency[failure]) + "        "
			failureData := []string{string(filePath), failure, numberOfCommits}
			data = append(data, failureData)
		}
	}
	return data
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
