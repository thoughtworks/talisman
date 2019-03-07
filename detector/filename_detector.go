package detector

import (
	"fmt"
	"regexp"

	"talisman/git_repo"

	log "github.com/Sirupsen/logrus"
)

//FileNameDetector represents tests performed against the fileName of the Additions.
//The Paths of the supplied Additions are tested against the configured patterns and if any of them match, it is logged as a failure during the run
type FileNameDetector struct {
	flagPatterns []*regexp.Regexp
}

//DefaultFileNameDetector returns a FileNameDetector that tests Additions against the pre-configured patterns
func DefaultFileNameDetector() Detector {
	return NewFileNameDetector("^.+_rsa$",
		"^.+_dsa$",
		"^.+_ed25519$",
		"^.+_ecdsa$",
		"^\\.\\w+_history$",
		"^.+\\.pem$",
		"^.+\\.ppk$",
		"^.+\\.key(pair)?$",
		"^.+\\.pkcs12$",
		"^.+\\.pfx$",
		"^.+\\.p12$",
		"^.+\\.asc$",
		"^\\.?htpasswd$",
		"^\\.?netrc$",
		"^.*\\.tblk$",
		"^.*\\.ovpn$",
		"^.*\\.kdb$",
		"^.*\\.agilekeychain$",
		"^.*\\.keychain$",
		"^.*\\.key(store|ring)$",
		"^jenkins\\.plugins\\.publish_over_ssh\\.BapSshPublisherPlugin.xml$",
		"^credentials\\.xml$",
		"^.*\\.pubxml(\\.user)?$",
		"^\\.?s3cfg$",
		"^\\.gitrobrc$",
		"^\\.?(bash|zsh)rc$",
		"^\\.?(bash_|zsh_)?profile$",
		"^\\.?(bash_|zsh_)?aliases$",
		"^secret_token.rb$",
		"^omniauth.rb$",
		"^carrierwave.rb$",
		"^schema.rb$",
		"^database.yml$",
		"^settings.py$",
		"^.*(config)(\\.inc)?\\.php$",
		"^LocalSettings.php$",
		"\\.?env",
		"\\bdump|dump\\b",
		"\\bsql|sql\\b",
		"\\bdump|dump\\b",
		"password",
		"backup",
		"private.*key",
		"(oauth).*(token)",
		"^.*\\.log$",
		"^\\.?kwallet$",
		"^\\.?gnucash$")
}

//NewFileNameDetector returns a FileNameDetector that tests Additions against the supplied patterns
func NewFileNameDetector(patternStrings ...string) Detector {
	var patterns = make([]*regexp.Regexp, len(patternStrings))
	for i, p := range patternStrings {
		patterns[i], _ = regexp.Compile(p)
	}
	return FileNameDetector{patterns}
}

//Test tests the fileNames of the Additions to ensure that they don't look suspicious
func (fd FileNameDetector) Test(additions []git_repo.Addition, ignoreConfig TalismanRCIgnore, result *DetectionResults) {
	cc := NewChecksumCompare(additions, ignoreConfig)
	for _, addition := range additions {
		if ignoreConfig.Deny(addition, "filename") || cc.IsScanNotRequired(addition){
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
				result.Fail(addition.Path, fmt.Sprintf("The file name %q failed checks against the pattern %s", addition.Path, pattern), addition.Commits)
			}
		}
	}
}
