package detector

import (
	"fmt"
	"regexp"
	"strings"

	"talisman/git_repo"

	log "github.com/Sirupsen/logrus"
)

type fn func(fc *FileContentDetector, word string) string

type FileContentDetector struct {
	base64Detector     *Base64Detector
	hexDetector        *HexDetector
	creditCardDetector *CreditCardDetector
}

func NewFileContentDetector() *FileContentDetector {
	fc := FileContentDetector{}
	fc.base64Detector = NewBase64Detector()
	fc.hexDetector = NewHexDetector()
	fc.creditCardDetector = NewCreditCardDetector()
	return &fc
}

func (fc *FileContentDetector) AggressiveMode() *FileContentDetector {
	fc.base64Detector.aggressiveDetector = &Base64AggressiveDetector{}
	return fc
}

func (fc *FileContentDetector) Test(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	cc := NewChecksumCompare(additions, ignoreConfig)
	for _, addition := range additions {
		if ignoreConfig.Deny(addition, "filecontent") || cc.IsScanNotRequired(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filecontent")
			continue
		}

		if string(addition.Name) == DefaultRCFileName {
			re := regexp.MustCompile(`(?i)checksum[ \t]*:[ \t]*[0-9a-fA-F]+`)
			content := re.ReplaceAllString(string(addition.Data), "")
			data := []byte(content)
			addition.Data = data
		}

		base64Results := fc.detectFile(addition.Data, checkBase64)
		fillBase46DetectionResults(base64Results, addition, result)

		hexResults := fc.detectFile(addition.Data, checkHex)
		fillHexDetectionResults(hexResults, addition, result)

		creditCardResults := fc.detectFile(addition.Data, checkCreditCardNumber)
		fillCreditCardDetectionResults(creditCardResults, addition, result)
	}
}

func fillResults(results []string, addition git_repo.Addition, result *DetectionResults, info string, output string) {
	for _, res := range results {
		if res != "" {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info(info)
			if string(addition.Name) == DefaultRCFileName {
				result.Warn(addition.Path, fmt.Sprintf(output, res), []string{})
			} else {
				result.Fail(addition.Path, fmt.Sprintf(output, res), []string{})
			}
		}
	}
}

func fillBase46DetectionResults(base64Results []string, addition git_repo.Addition, result *DetectionResults) {
	const info = "Failing file as it contains a base64 encoded text."
	const output = "Expected file to not to contain base64 encoded texts such as: %s"
	fillResults(base64Results, addition, result, info, output)
}

func fillCreditCardDetectionResults(creditCardResults []string, addition git_repo.Addition, result *DetectionResults) {
	const info = "Failing file as it contains a potential credit card number."
	const output = "Expected file to not to contain credit card numbers such as: %s"
	fillResults(creditCardResults, addition, result, info, output)
}

func fillHexDetectionResults(hexResults []string, addition git_repo.Addition, result *DetectionResults) {
	const info = "Failing file as it contains a hex encoded text."
	const output = "Expected file to not to contain hex encoded texts such as: %s"
	fillResults(hexResults, addition, result, info, output)
}

func (fc *FileContentDetector) detectFile(data []byte, getResult fn) []string {
	content := string(data)
	return fc.checkEachLine(content, getResult)
}

func (fc *FileContentDetector) checkEachLine(content string, getResult fn) []string {
	lines := strings.Split(content, "\n")
	res := []string{}
	for _, line := range lines {
		lineResult := fc.checkEachWord(line, getResult)
		if len(lineResult) > 0 {
			res = append(res, lineResult...)
		}
	}
	return res
}

func (fc *FileContentDetector) checkEachWord(line string, getResult fn) []string {
	words := strings.Fields(line)
	res := []string{}
	for _, word := range words {
		wordResult := getResult(fc, word)
		if wordResult != "" {
			res = append(res, wordResult)
		}
	}
	return res
}

func checkBase64(fc *FileContentDetector, word string) string {
	return fc.base64Detector.checkBase64Encoding(word)
}

func checkCreditCardNumber(fc *FileContentDetector, word string) string {
	return fc.creditCardDetector.checkCreditCardNumber(word)
}

func checkHex(fc *FileContentDetector, word string) string {
	return fc.hexDetector.checkHexEncoding(word)
}
