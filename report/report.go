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
func GenerateReport(r *detector.DetectionResults, directory string) string {

	var path string
	var jsonFilePath string
	var homeDir string
	var baseReportDirPath string

	usr, err := user.Current()
	homeDir = usr.HomeDir

	if directory == htmlReportDir {
		path = directory
		baseReportDirPath = filepath.Join(homeDir, ".talisman", htmlReportDir)
		jsonFilePath = filepath.Join(path, "/data", jsonFileName)
		os.MkdirAll(path, 0755)
		err = utility.Dir(baseReportDirPath, htmlReportDir)
		if err != nil {
			generateErrorMsg()
		}
	} else {
		path = filepath.Join(directory, "talisman_reports", "/data")
		jsonFilePath = filepath.Join(path, jsonFileName)
	}

	os.MkdirAll(path, 0755)

	jsonFile, err := os.Create(jsonFilePath)

	if err != nil {
		fmt.Printf("\n")
		log.Fatal("Cannot create report.json file\n", err)
	}

	jsonString, err := json.Marshal(r)
	if err != nil {
		log.Fatal("Unable to marshal JSON")
	}
	jsonFile.Write(jsonString)
	jsonFile.Close()
	return path
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
	fmt.Printf("\n")
}
