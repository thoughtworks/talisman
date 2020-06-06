package detector

import (
	"fmt"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

type FileSizeDetector struct {
	size int
}

func NewFileSizeDetector(size int) Detector {
	return FileSizeDetector{size}
}

func (fd FileSizeDetector) Test(comparator ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults) {
	for _, addition := range currentAdditions {
		if ignoreConfig.Deny(addition, "filesize") || comparator.IsScanNotRequired(addition) {
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
			result.Fail(addition.Path, "filesize", fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size), addition.Commits)
		}
	}
}
