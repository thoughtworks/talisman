package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"talisman/detector"
	"talisman/utility"

	"github.com/fatih/color"
)

const jsonFileName string = "report.json"
const htmlReportDir string = "talisman_html_report"

// GenerateReport generates a talisman scan report in html format
func GenerateReport(r *detector.DetectionResults, directory string) (path string, err error) {

	var jsonFilePath string
	var homeDir string
	var baseReportDirPath string

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %v", err.Error())
	}
	homeDir = usr.HomeDir

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
		path = filepath.Join(directory, "talisman_reports", "/data")
		jsonFilePath = filepath.Join(path, jsonFileName)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating path %s: %v", path, err)
	}

	jsonFile, err := os.Create(jsonFilePath)
	defer func() {
		if err = jsonFile.Close(); err != nil {
			err = fmt.Errorf("error closing file %s: %v", jsonFilePath, err)
		}
	}()

	if err != nil {
		return "", fmt.Errorf("error creating file %s: %v", jsonFilePath, err)
	}

	jsonString, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("error while marshal the report: %v", err)
	}
	_, err = jsonFile.Write(jsonString)
	if err != nil {
		return "", fmt.Errorf("error while writing report to file: %v", err)
	}
	return path, nil
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
