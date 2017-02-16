package detector

import (
	"encoding/base64"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
	"strings"
)

var delimiters = []string{".", "-", "="}

type FileContentDetector struct {
}

func NewFileContentDetector() Detector {
	return FileContentDetector{}
}

func (contentDetector FileContentDetector) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, addition := range additions {
		if ignores.Deny(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, fmt.Sprintf("%s was ignored by .talismanignore", addition.Path))
			continue
		}
		base64 := checkBase64EncodingForFile(addition.Data)
		if base64 == true {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Failing file as it contains a base64 encoded text.")
			result.Fail(addition.Path, fmt.Sprint("Expected file to not to contain base64 encoded texts"))
		}
	}
}

func checkBase64Encoding(s string) bool {
	if len(s) <= 4 {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func checkBase64EncodingForFile(content []byte) bool {
	s := string(content)
	for _, d := range delimiters {
		subStrings := strings.Split(s, d)
		if checkEachSubString(subStrings) {
			return true
		}
	}
	return false
}

func checkEachSubString(subStrings []string) bool {
	for _, sub := range subStrings {
		if checkBase64Encoding(sub) {
			return true
		}
	}
	return false
}
