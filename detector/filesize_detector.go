package detector

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
)

type FileSizeDetector struct {
	size int
}

func DefaultFileSizeDetector() Detector {
	return NewFileSizeDetector(1 * 1024 * 1024)
}

func NewFileSizeDetector(size int) Detector {
	return FileSizeDetector{size}
}

func (fd FileSizeDetector) Test(additions []git_repo.Addition, ignores Ignores, result *DetectionResults) {
	for _, addition := range additions {
		if ignores.Deny(addition, "filesize") {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filesize")
			continue
		}
		size := len(addition.Data)
		if size > fd.size {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
				"fileSize": size,
				"maxSize":  fd.size,
			}).Info("Failing file as it is larger than max allowed file size.")
			result.Fail(addition.Path, fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size))
		}
	}
}
