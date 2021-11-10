package main

import (
	"talisman/gitrepo"
	"talisman/utility"

	log "github.com/sirupsen/logrus"

	"github.com/bmatcuk/doublestar"
)

type PatternCmd struct {
	*runner
}

func NewPatternCmd(pattern string) *PatternCmd {
	var additions []gitrepo.Addition

	files, _ := doublestar.Glob(pattern)
	for _, file := range files {
		log.Debugf("reading file %s", file)
		data, err := utility.SafeReadFile(file)

		if err != nil {
			log.Warnf("Error reading file: %s. Skipping", file)
			continue
		}

		newAddition := gitrepo.NewAddition(file, data)
		additions = append(additions, newAddition)
	}

	return &PatternCmd{NewRunner(additions, "pattern")}
}
