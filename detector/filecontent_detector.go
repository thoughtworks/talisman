package detector

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
	"strings"
)

type FileContentDetector struct {
	base64Detector *Base64Detector
	hexDetector *HexDetector
}

func NewFileContentDetector() *FileContentDetector {
	fc := FileContentDetector{}
	fc.base64Detector = NewBase64Detector()
	fc.hexDetector = NewHexDetector()
	return &fc
}

func (fc *FileContentDetector) AggressiveMode() *FileContentDetector {
	fc.base64Detector.aggressiveDetector = &Base64AggressiveDetector{}
	return fc
}


func (fc *FileContentDetector) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, addition := range additions {
		if ignores.Deny(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, fmt.Sprintf("%s was ignored by .talismanignore", addition.Path))
			continue
		}
		base64Text := fc.detectFile(addition.Data)
		if base64Text != "" {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Failing file as it contains a base64 encoded text.")
			result.Fail(addition.Path, fmt.Sprintf("Expected file to not to contain base64 encoded texts such as: %s", base64Text))
		}
	}
}

func (fc *FileContentDetector) detectFile(data []byte) string {
	content := string(data)
	return fc.checkEachLine(content)
}

func (fc *FileContentDetector) checkEachLine(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		res := fc.checkEachWord(line)
		if res != "" {
			return res
		}
	}
	return ""
}

func (fc *FileContentDetector) checkEachWord(line string) string {
	words := strings.Fields(line)
	for _, word := range words {
		res := fc.base64Detector.checkBase64Encoding(word)
		if res != "" {
			return res
		}
		res = fc.hexDetector.checkHexEncoding(word)
		if res != "" {
			return res
		}
	}
	return ""
}