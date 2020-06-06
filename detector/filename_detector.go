package detector

import (
	"fmt"
	"regexp"

	"talisman/gitrepo"
	"talisman/talismanrc"

	log "github.com/Sirupsen/logrus"
)

var (
	filenamePatterns = []*regexp.Regexp{
		regexp.MustCompile(`^.+_rsa$`),
		regexp.MustCompile(`^.+_dsa.*$`),
		regexp.MustCompile(`^.+_ed25519$`),
		regexp.MustCompile(`^.+_ecdsa$`),
		regexp.MustCompile(`^\.\w+_history$`),
		regexp.MustCompile(`^.+\.pem$`),
		regexp.MustCompile(`^.+\.ppk$`),
		regexp.MustCompile(`^.+\.key(pair)?$`),
		regexp.MustCompile(`^.+\.pkcs12$`),
		regexp.MustCompile(`^.+\.pfx$`),
		regexp.MustCompile(`^.+\.p12$`),
		regexp.MustCompile(`^.+\.asc$`),
		regexp.MustCompile(`^\.?htpasswd$`),
		regexp.MustCompile(`^\.?netrc$`),
		regexp.MustCompile(`^.*\.tblk$`),
		regexp.MustCompile(`^.*\.ovpn$`),
		regexp.MustCompile(`^.*\.kdb$`),
		regexp.MustCompile(`^.*\.agilekeychain$`),
		regexp.MustCompile(`^.*\.keychain$`),
		regexp.MustCompile(`^.*\.key(store|ring)$`),
		regexp.MustCompile(`^jenkins\.plugins\.publish_over_ssh\.BapSshPublisherPlugin.xml$`),
		regexp.MustCompile(`^credentials\.xml$`),
		regexp.MustCompile(`^.*\.pubxml(\.user)?$`),
		regexp.MustCompile(`^\.?s3cfg$`),
		regexp.MustCompile(`^\.gitrobrc$`),
		regexp.MustCompile(`^\.?(bash|zsh)rc$`),
		regexp.MustCompile(`^\.?(bash_|zsh_)?profile$`),
		regexp.MustCompile(`^\.?(bash_|zsh_)?aliases$`),
		regexp.MustCompile(`^secret_token.rb$`),
		regexp.MustCompile(`^omniauth.rb$`),
		regexp.MustCompile(`^carrierwave.rb$`),
		regexp.MustCompile(`^schema.rb$`),
		regexp.MustCompile(`^database.yml$`),
		regexp.MustCompile(`^settings.py$`),
		regexp.MustCompile(`^.*(config)(\.inc)?\.php$`),
		regexp.MustCompile(`^LocalSettings.php$`),
		regexp.MustCompile(`\.?env`),
		regexp.MustCompile(`\bdump|dump\b`),
		regexp.MustCompile(`\bsql|sql\b`),
		regexp.MustCompile(`\bdump|dump\b`),
		regexp.MustCompile(`password`),
		regexp.MustCompile(`backup`),
		regexp.MustCompile(`private.*key`),
		regexp.MustCompile(`(oauth).*(token)`),
		regexp.MustCompile(`^.*\.log$`),
		regexp.MustCompile(`^\.?kwallet$`),
		regexp.MustCompile(`^\.?gnucash$`),
	}
)

//FileNameDetector represents tests performed against the fileName of the Additions.
//The Paths of the supplied Additions are tested against the configured patterns and if any of them match, it is logged as a failure during the run
type FileNameDetector struct {
	flagPatterns []*regexp.Regexp
}

//DefaultFileNameDetector returns a FileNameDetector that tests Additions against the pre-configured patterns
func DefaultFileNameDetector() Detector {
	return NewFileNameDetector(filenamePatterns)
}

//NewFileNameDetector returns a FileNameDetector that tests Additions against the supplied patterns
func NewFileNameDetector(patterns []*regexp.Regexp) Detector {
	return FileNameDetector{patterns}
}

//Test tests the fileNames of the Additions to ensure that they don't look suspicious
func (fd FileNameDetector) Test(allAdditions []gitrepo.Addition, currentAdditions []gitrepo.Addition, ignoreConfig *talismanrc.TalismanRC, result *DetectionResults) {
	cc := NewChecksumCompare(allAdditions, currentAdditions, ignoreConfig)
	for _, addition := range currentAdditions {
		if ignoreConfig.Deny(addition, "filename") || cc.IsScanNotRequired(addition) {
			log.WithFields(log.Fields{
				"filePath": addition.Path,
			}).Info("Ignoring addition as it was specified to be ignored.")
			result.Ignore(addition.Path, "filename")
			continue
		}
		for _, pattern := range fd.flagPatterns {
			if pattern.MatchString(string(addition.Name)) {
				log.WithFields(log.Fields{
					"filePath": addition.Path,
					"pattern":  pattern,
				}).Info("Failing file as it matched pattern.")
				result.Fail(addition.Path, "filename", fmt.Sprintf("The file name %q failed checks against the pattern %s", addition.Path, pattern), addition.Commits)
			}
		}
	}
}
