package detector

import (
	"regexp"
	"strings"

	"github.com/thoughtworks/talisman/git_repo"
)

const (
	//CommentPattern represents the prefix of a comment line in the ignore file
	CommentPattern string = "#"
	LinePattern string = "^([^#]+)?\\s*(#(.*))?$"

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
	return Ignore{pattern: pattern, comment: comment}
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
	return len(i.effectiveRules()) == 0
}

//Accept answers true if the Addition.Path is configured to be checked by the detectors
func (i Ignores) Accept(addition git_repo.Addition) bool {
	return !i.Deny(addition)
}

//Deny answers true if the Addition.Path is configured to be ignored and not checked by the detectors
func (i Ignores) Deny(addition git_repo.Addition) bool {
	result := false
	for _, pattern := range i.effectiveRules() {
		result = result || addition.Matches(pattern)
	}
	return result
}

func (i Ignores) effectiveRules() []string {
	var result []string
	for _, pattern := range i.patterns {
		if !isEmptyString(pattern.pattern) {
			result = append(result, pattern.pattern)
		}
	}
	return result
}

func isEmptyString(str string) bool {
	var emptyStringPattern = regexp.MustCompile("^\\s*$")
	return emptyStringPattern.MatchString(str)
}
