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

func (fd FileSizeDetector) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults, additionCompletionCallback func()) {
	severity := severity.SeverityConfiguration["LargeFileSize"]
	for _, addition := range currentAdditions {
		if ignoreConfig.Deny(addition, "filesize") || comparator.IsScanNotRequired(addition) {
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
			if severity.ExceedsThreshold(ignoreConfig.Threshold) {
				result.Fail(addition.Path, "filesize", fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size), addition.Commits, severity)
			} else {
				result.Warn(addition.Path, "filesize", fmt.Sprintf("The file name %q with file size %d is larger than max allowed file size(%d)", addition.Path, size, fd.size), addition.Commits, severity)
			}
		}
		additionCompletionCallback()
	}
}
