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
		if ignores.Deny(addition, "filecontent") {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filecontent")
			continue
		}
		base64Results := fc.detectFile(addition.Data)
		fillDetectionResults(base64Results, addition, result)
	}
}

func fillDetectionResults(base64Results []string, addition git_repo.Addition, result *DetectionResults) {
	for _, base64Res := range base64Results {
		if base64Res != "" {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Failing file as it contains a base64 encoded text.")
			result.Fail(addition.Path, fmt.Sprintf("Expected file to not to contain base64 or hex encoded texts such as: %s", base64Res))
		}
	}
}

func (fc *FileContentDetector) detectFile(data []byte) []string {
	content := string(data)
	return fc.checkEachLine(content)
}

func (fc *FileContentDetector) checkEachLine(content string) []string {
	lines := strings.Split(content, "\n")
	res := []string{}
	for _, line := range lines {
		lineResult := fc.checkEachWord(line)
		if len(lineResult) > 0 {
			res = append(res, lineResult...)
		}
	}
	return res
}

func (fc *FileContentDetector) checkEachWord(line string) []string {
	words := strings.Fields(line)
	res := []string{}
	for _, word := range words {
		wordResult := fc.base64Detector.checkBase64Encoding(word)
		if wordResult != "" {
			res = append(res, wordResult)
		}
		wordResult = fc.hexDetector.checkHexEncoding(word)
		if wordResult != "" {
			res = append(res, wordResult)
		}
	}
	return res
}
