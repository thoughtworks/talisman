package detector

import (
	"fmt"
	"os"
	"strings"
	"talisman/git_repo"
	"talisman/utility"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

type Details struct {
	Category string `json:"type"`
	Message string `json:"message"`
	Commits []string `json:"commits"`
}

type ResultsDetails struct {
	Filename git_repo.FilePath `json:"filename"`
	FailureList []Details      `json:"failure_list"`
	WarningList []Details      `json:"warning_list"`
	IgnoreList []Details 	   `json:"ignore_list"`
}

type FailureTypes struct  {
	Filecontent int `json:"filecontent"`
	Filesize int `json:"filesize"`
	Filename int `json:"filename"`
	Warnings int `json:"warnings"`
	Ignores int `json:"ignores"`
}

type ResultsSummary struct {
	Types FailureTypes `json:"types"`
}
//
//
//type FailureData struct {
//	FailuresInCommits map[string][]string
//}

//DetectionResults represents all interesting information collected during a detection run.
//It serves as a collecting parameter for the tests performed by the various Detectors in the DetectorChain
//Currently, it keeps track of failures and ignored files.
//The results are grouped by FilePath for easy reporting of all detected problems with individual files.
type DetectionResults struct {
	Summary ResultsSummary `json:"summary"`
	Results []ResultsDetails `json:"results"`
}

func (r *ResultsDetails) getWarningDataByCategoryAndMessage(failureMessage string, category string) *Details {
	detail := getDetaisByCategoryAndMessage(r.WarningList, category, failureMessage)
	r.WarningList = append(r.WarningList, *detail)
	return detail
}

func (r *ResultsDetails) getFailureDataByCategoryAndMessage(failureMessage string, category string) *Details {
	detail := getDetaisByCategoryAndMessage(r.FailureList, category, failureMessage)
	if detail == nil {
		detail = &Details{category, failureMessage, make([]string, 0)}
		r.FailureList = append(r.FailureList, *detail)
	}
	return detail
}

func (r *ResultsDetails) addIgnoreDataByCategory(category string) {
	isCategoryAlreadyPresent := false
	for _, detail := range r.IgnoreList {
		if strings.Compare(detail.Category, category) == 0 {
			isCategoryAlreadyPresent = true
		}
	}
	if !isCategoryAlreadyPresent {
		detail := Details{category, "", make([]string, 0)}
		r.IgnoreList = append(r.IgnoreList, detail)
	}
}

func getDetaisByCategoryAndMessage(detailsList []Details, category string, failureMessage string) *Details {
	for _, detail := range detailsList {
		if strings.Compare(detail.Category, category) == 0 && strings.Compare(detail.Message, failureMessage) == 0 {
			return &detail
		}
	}

	return nil
}

func (r *DetectionResults) getResultDetailsForFilePath(fileName git_repo.FilePath) *ResultsDetails {
	for _, resultDetail := range r.Results {
		if resultDetail.Filename == fileName {
			return &resultDetail
		}
 	}
	//resultDetail := ResultsDetails{fileName, make([]Details, 0), make([]Details, 0), make([]Details, 0)}
	//r.Results = append(r.Results, resultDetail)
	return nil
}

//NewDetectionResults is a new DetectionResults struct. It represents the pre-run state of a Detection run.
func NewDetectionResults() *DetectionResults {
	result := DetectionResults{ResultsSummary{FailureTypes{0,0,0, 0, 0}},make([]ResultsDetails, 0)}
	return &result
}


//Fail is used to mark the supplied FilePath as failing a detection for a supplied reason.
//Detectors are encouraged to provide context sensitive messages so that fixing the errors is made simple for the end user
//Fail may be called multiple times for each FilePath and the calls accumulate the provided reasons
func (r *DetectionResults) Fail(filePath git_repo.FilePath, category string, message string, commits []string) {
	isFilePresentInResults := false
	for resultIndex := 0; resultIndex < len(r.Results); resultIndex++ {
		if r.Results[resultIndex].Filename == filePath {
			isFilePresentInResults = true
			isEntryPresentForGivenCategoryAndMessage := false
			for detailIndex := 0; detailIndex < len(r.Results[resultIndex].FailureList); detailIndex++ {
				if strings.Compare(r.Results[resultIndex].FailureList[detailIndex].Category, category) == 0 && strings.Compare(r.Results[resultIndex].FailureList[detailIndex].Message, message) == 0 {
					isEntryPresentForGivenCategoryAndMessage = true
					r.Results[resultIndex].FailureList[detailIndex].Commits = append(r.Results[resultIndex].FailureList[detailIndex].Commits, commits...)
				}
			}
			if !isEntryPresentForGivenCategoryAndMessage {
				r.Results[resultIndex].FailureList = append(r.Results[resultIndex].FailureList, Details{category, message, commits})
			}
		}
	}
	if !isFilePresentInResults {
		failureDetails := Details{category, message, commits}
		resultDetails := ResultsDetails{filePath, make([]Details, 0), make([]Details, 0), make([]Details, 0)}
		resultDetails.FailureList = append(resultDetails.FailureList, failureDetails)
		r.Results = append(r.Results, resultDetails)
	}
	r.updateResultsSummary(category)
}

func (r *DetectionResults) Warn(filePath git_repo.FilePath, category string, message string, commits []string) {
	isFilePresentInResults := false
	for resultIndex := 0; resultIndex < len(r.Results); resultIndex++ {
		if r.Results[resultIndex].Filename == filePath {
			isFilePresentInResults = true
			isEntryPresentForGivenCategoryAndMessage := false
			for detailIndex := 0; detailIndex < len(r.Results[resultIndex].WarningList); detailIndex++ {
				if strings.Compare(r.Results[resultIndex].WarningList[detailIndex].Category, category) == 0 && strings.Compare(r.Results[resultIndex].WarningList[detailIndex].Message, message) == 0 {
					isEntryPresentForGivenCategoryAndMessage = true
					r.Results[resultIndex].WarningList[detailIndex].Commits = append(r.Results[resultIndex].WarningList[detailIndex].Commits, commits...)
				}
			}
			if !isEntryPresentForGivenCategoryAndMessage {
				r.Results[resultIndex].WarningList = append(r.Results[resultIndex].WarningList, Details{category, message, commits})
			}
		}
	}
	if !isFilePresentInResults {
		warningDetails := Details{category, message, commits}
		resultDetails := ResultsDetails{filePath, make([]Details, 0), make([]Details, 0), make([]Details, 0)}
		resultDetails.WarningList = append(resultDetails.WarningList, warningDetails)
		r.Results = append(r.Results, resultDetails)
	}
	r.Summary.Types.Warnings++
}

//Ignore is used to mark the supplied FilePath as being ignored.
//The most common reason for this is that the FilePath is Denied by the Ignores supplied to the Detector, however, Detectors may use more sophisticated reasons to ignore files.
func (r *DetectionResults) Ignore(filePath git_repo.FilePath, category string) {

	isFilePresentInResults := false
	for resultIndex := 0; resultIndex < len(r.Results); resultIndex++ {
		if r.Results[resultIndex].Filename == filePath {
			isFilePresentInResults = true
			isEntryPresentForGivenCategory := false
			for detailIndex := 0; detailIndex < len(r.Results[resultIndex].IgnoreList); detailIndex++ {
				if strings.Compare(r.Results[resultIndex].IgnoreList[detailIndex].Category, category) == 0 {
					isEntryPresentForGivenCategory = true

				}
			}
			if !isEntryPresentForGivenCategory {
				detail := Details{category, "", make([]string, 0)}
				r.Results[resultIndex].IgnoreList = append(r.Results[resultIndex].IgnoreList, detail)
			}
		}
	}
	if !isFilePresentInResults {
		ignoreDetails := Details{category, "", make([]string, 0)}
		resultDetails := ResultsDetails{filePath, make([]Details, 0), make([]Details, 0), make([]Details, 0)}
		resultDetails.IgnoreList = append(resultDetails.IgnoreList, ignoreDetails)
		r.Results = append(r.Results, resultDetails)
	}
	r.Summary.Types.Ignores++
}


func createNewResultForFile(category string, message string, commits []string, filePath git_repo.FilePath) ResultsDetails {
	failureDetails := Details{category, message, commits}
	resultDetails := ResultsDetails{filePath, make([]Details, 0), make([]Details, 0), make([]Details, 0)}
	resultDetails.FailureList = append(resultDetails.FailureList, failureDetails)
	return resultDetails
}

func (r *DetectionResults) updateResultsSummary(category string) {
	if strings.Compare("filecontent", category) == 0 {
		r.Summary.Types.Filecontent++
	} else if strings.Compare("filename", category) == 0 {
		r.Summary.Types.Filename++
	} else if strings.Compare("filesize", category) == 0 {
		r.Summary.Types.Filesize++
	}

}

//HasFailures answers if any Failures were detected for any FilePath in the current run
func (r *DetectionResults) HasFailures() bool {
	return r.Summary.Types.Filesize > 0 || r.Summary.Types.Filename > 0 || r.Summary.Types.Filecontent > 0
}

//HasIgnores answers if any FilePaths were ignored in the current run
func (r *DetectionResults) HasIgnores() bool {
	return r.Summary.Types.Ignores > 0
}

func (r *DetectionResults) HasWarnings() bool {
	return r.Summary.Types.Warnings > 0
}

func (r *DetectionResults) HasDetectionMessages() bool {
	return r.HasWarnings() || r.HasFailures() || r.HasIgnores()
}

//Successful answers if no detector was able to find any possible result to fail the run
func (r *DetectionResults) Successful() bool {
	return !r.HasFailures()
}

//GetFailures returns the various reasons that a given FilePath was marked as failing by all the detectors in the current run
func (r *DetectionResults) GetFailures(fileName git_repo.FilePath) []Details {
	return r.getResultDetailsForFilePath(fileName).FailureList
}

func (r *DetectionResults) ReportWarnings() string {
	var result string
	var filePathsForWarnings []string
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Warnings"})
	table.SetRowLine(true)

	for _, resultDetails := range r.Results {
		if len(resultDetails.WarningList) > 0 {
			filePathsForWarnings = append(filePathsForWarnings, string(resultDetails.Filename))
			warningData := r.ReportFileWarnings(resultDetails.Filename)
			data = append(data, warningData...)
		}
	}


	filePathsForWarnings = utility.UniqueItems(filePathsForWarnings)
	if r.Summary.Types.Warnings > 0 {
		fmt.Printf("\n\x1b[1m\x1b[31mTalisman Warnings:\x1b[0m\x1b[0m\n")
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

	for _, resultDetails := range r.Results {
		if len(resultDetails.FailureList) > 0 || len(resultDetails.IgnoreList) > 0 {
			filePathsForIgnoresAndFailures = append(filePathsForIgnoresAndFailures, string(resultDetails.Filename))
			failureData := r.ReportFileFailures(resultDetails.Filename)
			data = append(data, failureData...)
		}
	}

	filePathsForIgnoresAndFailures = utility.UniqueItems(filePathsForIgnoresAndFailures)

	if r.HasFailures() {
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
	failureList := r.getResultDetailsForFilePath(filePath).FailureList
	var data [][]string
	if len(failureList) > 0 {
		for _, detail := range failureList {
			if len(detail.Message) > 150 {
				detail.Message = detail.Message[:150] + "\n" + detail.Message[150:]
			}
			data = append(data, []string{string(filePath), detail.Message})
		}
	}
	return data
}

func (r *DetectionResults) ReportFileWarnings(filePath git_repo.FilePath) [][]string {
	warningList := r.getResultDetailsForFilePath(filePath).WarningList
	var data [][]string
	if len(warningList) > 0 {
		for _, detail := range warningList {
			if len(detail.Message) > 150 {
				detail.Message = detail.Message[:150] + "\n" + detail.Message[150:]
			}
			data = append(data, []string{string(filePath), detail.Message})
		}
	}
	return data
}

func keys(aMap map[git_repo.FilePath][]string) []git_repo.FilePath {
	var result []git_repo.FilePath
	for filePath := range aMap {
		result = append(result, filePath)
	}
	return result
}

