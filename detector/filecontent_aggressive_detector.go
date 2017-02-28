package detector

import (
	"encoding/base64"
	"strings"
)

var delimiters = []string{".", "-", "="}
const aggressivenessThreshold = 15 //decreasing makes it more aggressive

type AggressiveFileContentDetector struct {
}

func (ac *AggressiveFileContentDetector) Test(s string) string {
	for _, d := range delimiters {
		subStrings := strings.Split(s, d)
		res := checkEachSubString(subStrings)
		if res != "" {
			return res
		}
	}
	return ""
}

func decodeBase64(s string) string {
	if len(s) <= aggressivenessThreshold {
		return ""
	}
	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return s
}

func checkEachSubString(subStrings []string) string {
	for _, sub := range subStrings {
		suspicious := decodeBase64(sub)
		if suspicious != "" {
			return suspicious
		}
	}
	return ""
}