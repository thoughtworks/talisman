package detector

import (
	"regexp"
	"strings"

	"github.com/thoughtworks/talisman/git_repo"
)

const (
	//LinePattern represents a line in the ignorefile with an optional comment
	LinePattern string = "^([^#]+)?\\s*(#(.*))?$"

	//IgnoreDetectorCommentPattern represents a special comment that ignores only certain detectors
	IgnoreDetectorCommentPattern string = "^ignore:([^\\s]+).*$"

	//DefaultIgnoreFileName represents the name of the default file in which the ignore patterns are configured
	DefaultIgnoreFileName string = ".talismanignore"
)

//Ignores represents a set of patterns that have been configured to be ignored by the Detectors.
//Detectors are expected to honor these ignores.
type Ignores struct {
	patterns []Ignore
}
//Ignore represents a single pattern and its comment
type Ignore struct {
	pattern string
	comment string
	ignoredDetectors []string
}

//ReadIgnoresFromFile builds an Ignores from the lines configured in a File.
//The file itself is supplied as a File Read operation, which is specified, by default, as reading a file in the root of the repository.
//The file name that is read is DEFAULT_IGNORE_FILE_NAME (".talismanignore")
func ReadIgnoresFromFile(repoFileRead func(string) ([]byte, error)) Ignores {
	contents, err := repoFileRead(DefaultIgnoreFileName)
	if err != nil {
		panic(err)
	}
	return NewIgnores(strings.Split(string(contents), "\n")...)
}

func NewIgnore(pattern string, comment string) Ignore {
	var ignoredDetectors []string
	ignorePattern := regexp.MustCompile(IgnoreDetectorCommentPattern)
	match := ignorePattern.FindStringSubmatch(comment)
	if match != nil {
		ignoredDetectors = strings.Split(match[1], ",")
	}

	return Ignore{
		pattern: pattern,
		comment: comment,
		ignoredDetectors: ignoredDetectors,
	}
}

func (i Ignore) isEffective(detectorName string) bool {
	return !isEmptyString(i.pattern) &&
		(contains(i.ignoredDetectors, detectorName) || len(i.ignoredDetectors) == 0)
}
//NewIgnores builds a new Ignores with the patterns specified in the ignoreSpecs
//Empty lines and comments are ignored.
func NewIgnores(lines ...string) Ignores {
	var ignores []Ignore
	for _, line := range lines {
		var commentPattern = regexp.MustCompile(LinePattern)
		groups := commentPattern.FindStringSubmatch(line)
		if len(groups) == 4 {
			ignores = append(ignores, NewIgnore(strings.TrimSpace(groups[1]), strings.TrimSpace(groups[3])))
		}
	}
	return Ignores{ignores}
}

//AcceptsAll returns true if there are no rules specified
func (i Ignores) AcceptsAll() bool {
	return len(i.effectiveRules("any-detector")) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (i Ignores) Accept(addition git_repo.Addition, detectorName string) bool {
	return !i.Deny(addition, detectorName)
}

//Deny answers true if the Addition.Path is configured to be ignored and not checked by the detectors
func (i Ignores) Deny(addition git_repo.Addition, detectorName string) bool {
	result := false
	for _, pattern := range i.effectiveRules(detectorName) {
		result = result || addition.Matches(pattern)
	}
	return result
}

func (i Ignores) effectiveRules(detectorName string) []string {
	var result []string
	for _, ignore := range i.patterns {
		if ignore.isEffective(detectorName) {
			result = append(result, ignore.pattern)
		}
	}
	return result
}

func isEmptyString(str string) bool {
	var emptyStringPattern = regexp.MustCompile("^\\s*$")
	return emptyStringPattern.MatchString(str)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
