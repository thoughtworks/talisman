package filesize

import (
	"fmt"
	"talisman/detector/detector"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/sirupsen/logrus"
)

type FileSizeDetector struct {
	size int
}

func NewFileSizeDetector(size int) detector.Detector {
	return FileSizeDetector{size}
}

func (fd FileSizeDetector) Test(comparator helpers.IgnoreEvaluator, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	largeFileSizeSeverity := severity.SeverityConfiguration["LargeFileSize"]
	for _, addition := range currentAdditions {
		if comparator.ShouldIgnore(addition, "filesize") {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filesize")
			additionCompletionCallback()
			continue
		}
		size := len(addition.Data)
		if size > fd.size {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
				"fileSize": size,
				"maxSize":  fd.size,
			}).Info("Failing file as it is larger than max allowed file size.")
			if largeFileSizeSeverity.ExceedsThreshold(ignoreConfig.Threshold) {
				result.Fail(addition.Path, "filesize", fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size), addition.Commits, largeFileSizeSeverity)
			} else {
				result.Warn(addition.Path, "filesize", fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size), addition.Commits, largeFileSizeSeverity)
			}
		}
		additionCompletionCallback()
	}
}
