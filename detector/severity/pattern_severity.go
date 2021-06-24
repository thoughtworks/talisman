package severity

import (
	"regexp"
)

type PatternSeverity struct {
	Pattern  *regexp.Regexp
	Severity Severity
}


