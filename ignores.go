package main

import (
	"regexp"
	"strings"

	"github.com/thoughtworks/talisman/git_repo"
)

const (
	//COMMENT_PATTERN represents the prefix that needs to be attached to a line in the .talismanignore file to mark it as a comment
	COMMENT_PATTERN string = "#"

	//DEFAULT_IGNORE_FILE_NAME represents the name of the default file in which the ignore patterns are configured
	DEFAULT_IGNORE_FILE_NAME string = ".talismanignore"
)

//Ignores represents a set of patterns that have been configured to be ignored by the Detectors.
//Detectors are expected to honor these ignores.
type Ignores struct {
	patterns []string
}

//ReadIgnoresFromFile builds an Ignores from the lines configured in a File.
//The file itself is supplied as a File Read operation, which is specified, by default, as reading a file in the root of the repository.
//The file name that is read is DEFAULT_IGNORE_FILE_NAME (".talismanignore")
func ReadIgnoresFromFile(repoFileRead func(string) ([]byte, error)) Ignores {
	contents, err := repoFileRead(DEFAULT_IGNORE_FILE_NAME)
	if err != nil {
		panic(err)
	}
	var trimmedLines []string
	for _, line := range strings.Split(string(contents), "\n") {
		trimmedLines = append(trimmedLines, strings.TrimSpace(line))
	}
	return NewIgnores(trimmedLines...)
}

//NewIgnores builds a new Ignores with the patterns specified in the ignoreSpecs
//Empty lines and comments are ignored.
func NewIgnores(ignoreSpecs ...string) Ignores {
	return Ignores{ignoreSpecs}
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
		if !isEmptyOrComment(pattern) {
			result = append(result, pattern)
		}
	}
	return result
}

func isEmptyOrComment(pattern string) bool {
	return isEmptyString(pattern) || strings.HasPrefix(pattern, COMMENT_PATTERN)
}

func isEmptyString(str string) bool {
	var emptyStringPattern = regexp.MustCompile("^\\s*$")
	return emptyStringPattern.MatchString(str)
}
