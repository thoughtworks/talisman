package severity

import (
	"regexp"
)

type SeverityValue int

const (
	LowSeverity    SeverityValue = 1
	MediumSeverity SeverityValue = 2
	HighSeverity   SeverityValue = 3
)

type PatternSeverity struct {
	Pattern  *regexp.Regexp
	Severity Severity
}

func SeverityDisplayString(severity SeverityValue) string {
	switch severity {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	default:
		return "Undefined"
	}
}

type Severity struct {
	Value SeverityValue
}

func (s Severity) String() string {
	return SeverityDisplayString(s.Value)
}
func (s Severity) ExceedsThreshold(threshold SeverityValue) bool {
	return s.Value >= threshold
}
func Low() Severity {
	return Severity{
		Value: LowSeverity,
	}
}

func Medium() Severity {
	return Severity{
		Value: MediumSeverity,
	}
}

func High() Severity {
	return Severity{
		Value: HighSeverity,
	}
}
