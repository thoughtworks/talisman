package report

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"talisman/detector"
)

var reportsFolder string
const htmlFileName string = "report.html"
const jsonFileName string = "report.json"


// GenerateReport generates a talisman scan report in html format
func GenerateReport(r *detector.DetectionResults, directory string) string {

	var path string
	var jsonFilePath string
	var home_dir string
	var base_report_dir_path string

	if directory == "talisman_html_report" {
			path = directory
			// base_report_dir_path = filepath.Join(home_dir, ".talisman/bin/talisman_html_report")
			// Dir(base_report_dir_path, "talisman_html_report")
	} else {
			path = filepath.Join(directory, "talisman_reports")
	}

	jsonFilePath = filepath.Join(path, "/data", jsonFileName)


	usr, err := user.Current()
	home_dir = usr.HomeDir

	base_report_dir_path = filepath.Join(home_dir, ".talisman/bin/talisman_html_report")

	os.MkdirAll(path, 0755)
	Dir(base_report_dir_path, "talisman_html_report")

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatal("Cannot create report.json file", err)

	}
	//jsonResultSchema := detector.GetJsonSchema(r)
	jsonString, err := json.Marshal(r)
	if err != nil {
		log.Fatal("Unable to marshal JSON")
	}
	jsonFile.Write(jsonString)
	jsonFile.Close()
	return path
}

func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func Dir(src string, dst string) error {

	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
