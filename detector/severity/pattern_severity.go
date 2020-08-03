package severity

import (
	"regexp"
)

type SeverityValue int

const (
	LowSeverity = SeverityValue(iota + 1)
	MediumSeverity
	HighSeverity
)

type PatternSeverity struct {
	Pattern  *regexp.Regexp
	Severity Severity
}

type Severity struct {
	Value SeverityValue
}

func (s Severity) String() string {
	return SeverityValueToString(s.Value)
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
