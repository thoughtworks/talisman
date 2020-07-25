package filename

import (
	"fmt"
	"regexp"
	"talisman/detector/detector"
	"talisman/detector/helpers"
	"talisman/detector/severity"

	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

var (
	filenamePatterns = []*severity.PatternSeverity{
		{Pattern: regexp.MustCompile(`^.+_rsa$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+_dsa.*$`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`^.+_ed25519$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+_ecdsa$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.\w+_history$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+\.pem$`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`^.+\.ppk$`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`^.+\.key(pair)?$`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`^.+\.pkcs12$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+\.pfx$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+\.p12$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.+\.asc$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?htpasswd$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?netrc$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.tblk$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.ovpn$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.kdb$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.agilekeychain$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.keychain$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.key(store|ring)$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^jenkins\.plugins\.publish_over_ssh\.BapSshPublisherPlugin.xml$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^credentials\.xml$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.pubxml(\.user)?$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?s3cfg$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.gitrobrc$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?(bash|zsh)rc$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?(bash_|zsh_)?profile$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?(bash_|zsh_)?aliases$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^secret_token.rb$`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`^omniauth.rb$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^carrierwave.rb$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^schema.rb$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^database.yml$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^settings.py$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*(config)(\.inc)?\.php$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^LocalSettings.php$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`\.?env`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`\bdump|dump\b`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`\bsql|sql\b`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`\bdump|dump\b`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`password`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`backup`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`private.*key`), Severity: severity.High()},
		{Pattern: regexp.MustCompile(`(oauth).*(token)`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^.*\.log$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?kwallet$`), Severity: severity.Low()},
		{Pattern: regexp.MustCompile(`^\.?gnucash$`), Severity: severity.Low()},
	}
)

//FileNameDetector represents tests performed against the fileName of the Additions.
//The Paths of the supplied Additions are tested against the configured patterns and if any of them match, it is logged as a failure during the run
type FileNameDetector struct {
	flagPatterns []*severity.PatternSeverity
	threshold    severity.SeverityValue
}

//DefaultFileNameDetector returns a FileNameDetector that tests Additions against the pre-configured patterns
func DefaultFileNameDetector(threshold severity.SeverityValue) detector.Detector {
	return NewFileNameDetector(filenamePatterns, threshold)
}

//NewFileNameDetector returns a FileNameDetector that tests Additions against the supplied patterns
func NewFileNameDetector(patternsWithSeverity []*severity.PatternSeverity, threshold severity.SeverityValue) detector.Detector {
	return FileNameDetector{patternsWithSeverity, threshold}
}

//Test tests the fileNames of the Additions to ensure that they don't look suspicious
func (fd FileNameDetector) Test(comparator helpers.ChecksumCompare, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *helpers.DetectionResults) {
	for _, addition := range currentAdditions {
		if ignoreConfig.Deny(addition, "filename") || comparator.IsScanNotRequired(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filename")
			continue
		}
		for _, patternWithSeverity := range fd.flagPatterns {
			if patternWithSeverity.Pattern.MatchString(string(addition.Name)) {
				log.WithFields(log.Fields{
					"filePath": addition.Path,
					"pattern":  patternWithSeverity.Pattern,
					"severity": patternWithSeverity.Severity,
				}).Info("Failing file as it matched pattern.")
				if patternWithSeverity.Severity.ExceedsThreshold(fd.threshold) {
					result.Fail(addition.Path, "filename", fmt.Sprintf("The file name %q failed checks against the pattern %s", addition.Path, patternWithSeverity.Pattern), addition.Commits, patternWithSeverity.Severity)
				} else {
					result.Warn(addition.Path, "filename", fmt.Sprintf("The file name %q failed checks against the pattern %s", addition.Path, patternWithSeverity.Pattern), addition.Commits, patternWithSeverity.Severity)
				}
			}
		}
	}
}
