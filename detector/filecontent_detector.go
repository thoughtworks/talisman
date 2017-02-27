package detector

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
	"strings"
	"math"
)

const BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

type FileContentDetector struct {
	base64Map map[string]bool
}

func NewFileContentDetector() Detector {
	fc := FileContentDetector{}
	fc.initBase64Map()
	return &fc
}

func (fc *FileContentDetector) initBase64Map() {
	fc.base64Map = map[string]bool{}
	for i := 0; i < len(BASE64_CHARS); i++ {
		fc.base64Map[string(BASE64_CHARS[i])] = true
	}
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
		base64 := fc.checkBase64EncodingForFile(addition.Data)
		if base64 == true {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Failing file as it contains a base64 encoded text.")
			result.Fail(addition.Path, fmt.Sprint("Expected file to not to contain base64 encoded texts"))
		}
	}
}

func (fc *FileContentDetector) getShannonEntropy(str string, superSet string) float64 {
	if str == "" {
		return 0
	}
	entropy := 0.0
	for _, c := range superSet {
		p := float64(strings.Count(str, string(c))) / float64(len(str))
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func (fc *FileContentDetector) checkBase64Encoding(word string) bool {
	entropyCandidates := fc.getEntropyCandidatesWithinWord(word, 20, fc.base64Map)
	for _, candidate := range entropyCandidates {
		entropy := fc.getShannonEntropy(candidate, BASE64_CHARS)
		if entropy > 4.5 {
			return true
		}
	}
	return false
}

func (fc *FileContentDetector) getEntropyCandidatesWithinWord(word string, threshold int, superSet map[string]bool) []string {
	candidates := []string{}
	count := 0
	subSet := ""
	for _, c := range word {
		char := string(c)
		if _, ok := superSet[char]; ok {
			subSet += char
			count++
		} else {
			if count > threshold {
				candidates = append(candidates, subSet)
			}
			subSet = ""
			count = 0
		}
	}
	if count > threshold {
		candidates = append(candidates, subSet)
	}
	return candidates

}

func (fc *FileContentDetector) checkBase64EncodingForFile(data []byte) bool {
	content := string(data)
	isBase64 := fc.checkEachLine(content)
	if isBase64 {
		return true
	}
	return false
}

func (fc *FileContentDetector) checkEachLine(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if fc.checkEachWord(line) {
			return true
		}
	}
	return false
}

func (fc *FileContentDetector) checkEachWord(line string) bool {
	words := strings.Fields(line)
	for _, word := range words {
		if fc.checkBase64Encoding(word) {
			return true
		}
	}
	return false
}
