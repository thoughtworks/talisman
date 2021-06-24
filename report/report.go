package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"talisman/detector/helpers"
	"talisman/utility"

	"github.com/fatih/color"
)

const jsonFileName string = "report.json"
const htmlReportDir string = "talisman_html_report"
const jsonReportDir string = "talisman_html_report"

// GenerateReport generates a talisman scan report in html format
func GenerateReport(r *helpers.DetectionResults, directory string) (string, error) {
	var jsonFilePath string
	var homeDir string
	var baseReportDirPath string

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %v", err.Error())
	}
	homeDir = usr.HomeDir
	path := jsonReportDir
	if directory == htmlReportDir {
		path = directory
		baseReportDirPath = filepath.Join(homeDir, ".talisman", htmlReportDir)
		jsonFilePath = filepath.Join(path, "/data", jsonFileName)
		err = utility.Dir(baseReportDirPath, htmlReportDir)
		if err != nil {
			generateErrorMsg()
			return "", fmt.Errorf("error copying reports: %v", err)
		}
	} else {
		path = filepath.Join(directory, "talisman_reports")
		_ = os.RemoveAll(path)
		path = filepath.Join(path, "data")
		jsonFilePath = filepath.Join(path, jsonFileName)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating path %s: %v", path, err)
	}

	_, err = generateAndWriteToFile(r, jsonFilePath)
	if err != nil {
		return "", err
	}
	return path, nil
}

func generateAndWriteToFile(r *helpers.DetectionResults, jsonFilePath string) (path string, err error) {
	jsonFile, err := os.Create(jsonFilePath)
	defer func() {
		if err = jsonFile.Close(); err != nil {
			err = fmt.Errorf("error closing file %s: %v %#v", jsonFilePath, err, err)
		}
	}()

	if err != nil {
		return "", fmt.Errorf("error creating file %s: %v", jsonFilePath, err)
	}

	jsonString, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("error while rendering report json: %v : %#v", err, err)
	}
	_, err = jsonFile.Write(jsonString)
	if err != nil {
		return "", fmt.Errorf("error while writing report json to file: %v %#v", err, err)
	}
	return jsonFilePath, nil
}

func generateErrorMsg() {
	color.HiMagenta("\nLooks like you are using 'talisman --scanWithHtml' for scanning.")
	color.HiMagenta("But it appears that you have not installed Talisman Html Report")
	color.HiMagenta("Please go through Talisman Readme and make sure you install the same from:")
	color.Yellow("\nhttps://github.com/jaydeepc/talisman-html-report")
	color.Cyan("\nOR use 'talisman --scan' if you want the JSON report alone\n")
	fmt.Printf("\n")
	color.Red("Failed: Unable to perform Scan")
	fmt.Printf("\n")
	log.Fatalln("Run Status: Failed")
}
