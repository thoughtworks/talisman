package detector

import "talisman/git_repo"

type Details struct {
	Category string `json:"type"`
	Message string `json:"message"`
	Commits []string `json:"commits"`
}

type ResultsDetails struct {
	Filename git_repo.FilePath `json:"filename"`
	FailureList []Details      `json:"failure_list"`
	WarningList []Details      `json:"warning_list"`
}

type FailureTypes struct  {
	Filecontent int `json:"filecontent"`
	Filesize int `json:"filesize"`
	Filename int `json:"filename"`
}

type ResultsSummary struct {
	Types FailureTypes `json:"types"`
}

type JsonDetectionResults struct {
	Summary ResultsSummary `json:"summary"`
	Results []ResultsDetails `json:"results"`

}

func (result JsonDetectionResults) getResultObjectForFileName(filename git_repo.FilePath) ResultsDetails {
	for _, resultDetails := range result.Results {
		if resultDetails.Filename == filename {
			return resultDetails
		}
	}
	return ResultsDetails{"", make([]Details, 0), make([]Details, 0)}
}

func GetJsonSchema(r *DetectionResults) JsonDetectionResults {

	jsonResults := JsonDetectionResults{}
	failures := r.Failures
	for path, data := range failures {
		resultDetails := ResultsDetails{}
		resultDetails.Filename = path

		for message, commits := range data.FailuresInCommits {
			failureDetails := Details{}
			failureDetails.Message = message
			failureDetails.Commits = commits
			resultDetails.FailureList = append(resultDetails.FailureList, failureDetails)
		}
		jsonResults.Results = append(jsonResults.Results, resultDetails)
	}
	warnings := r.warnings
	for path, data := range warnings {
		resultDetails := jsonResults.getResultObjectForFileName(path)
		resultDetails.Filename = path

		for message, commits := range data.FailuresInCommits {
			failureDetails := Details{}
			failureDetails.Message = message
			failureDetails.Commits = commits
			resultDetails.WarningList = append(resultDetails.WarningList, failureDetails)
		}
	}
	return jsonResults
}

