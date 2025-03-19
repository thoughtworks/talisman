package detector

import (
	"talisman/detector/helpers"
	"talisman/gitrepo"
	"talisman/talismanrc"
)

// Detector represents a single kind of test to be performed against a set of Additions
// Detectors are expected to honor the ignores that are passed in and log them in the results
// Detectors are expected to signal any errors to the results
type Detector interface {
	Test(comparator helpers.IgnoreEvaluator, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func())
}
