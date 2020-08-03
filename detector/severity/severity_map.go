package severity

import "strings"

var severityMap = map[SeverityValue]string{
	LowSeverity:    "low",
	MediumSeverity: "medium",
	HighSeverity:   "high",
}

func SeverityValueToString(severity SeverityValue) string {
	return severityMap[severity]
}

func SeverityStringToValue(severity string) SeverityValue {
	severityInLowerCase := strings.ToLower(severity)
	for k, v := range severityMap {
		if v == severityInLowerCase {
			return k
		}
	}
	return 0
}
