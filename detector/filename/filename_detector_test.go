package filename

//This is completely derived from the really useful work done by Jen Andre here:
//https://github.com/jandre/safe-commit-hook

import (
	"regexp"
	"talisman/detector/helpers"
	"talisman/detector/severity"
	"talisman/utility"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

var talismanRC = &talismanrc.TalismanRC{}

func TestShouldFlagPotentialSSHPrivateKeys(t *testing.T) {
	shouldFail("id_rsa", "^.+_rsa$", severity.Low, t)
	shouldFail("id_dsa", "^.+_dsa.*$", severity.Low, t)
	shouldFail("id_dsa.pub", "^.+_dsa.*$", severity.Low, t)
	shouldFail("id_ed25519", "^.+_ed25519$", severity.Low, t)
	shouldFail("id_ecdsa", "^.+_ecdsa$", severity.Low, t)
}

func TestShouldFlagPotentialHistoryFiles(t *testing.T) {
	shouldFail(".bash_history", "^\\.\\w+_history$", severity.Low, t)
	shouldFail(".zsh_history", "^\\.\\w+_history$", severity.Low, t)
	shouldFail(".z_history", "^\\.\\w+_history$", severity.Low, t)
	shouldFail(".irb_history", "^\\.\\w+_history$", severity.Low, t)
	shouldFail(".psql_history", "^\\.\\w+_history$", severity.Low, t)
	shouldFail(".mysql_history", "^\\.\\w+_history$", severity.Low, t)
}

func TestShouldFlagPotentialPrivateKeys(t *testing.T) {
	shouldFail("foo.pem", "^.+\\.pem$", severity.Low, t)
	shouldFail("foo.ppk", "^.+\\.ppk$", severity.Low, t)
	shouldFail("foo.key", "^.+\\.key(pair)?$", severity.Low, t)
	shouldFail("foo.keypair", "^.+\\.key(pair)?$", severity.Low, t)
}

func TestShouldFlagPotentialKeyBundles(t *testing.T) {
	shouldFail("foo.pkcs12", "^.+\\.pkcs12$", severity.Low, t)
	shouldFail("foo.pfx", "^.+\\.pfx$", severity.Low, t)
	shouldFail("foo.p12", "^.+\\.p12$", severity.Low, t)
	shouldFail("foo.asc", "^.+\\.asc$", severity.Low, t)
}

func TestShouldFlagPotentialConfigurationFiles(t *testing.T) {
	shouldFail(".htpasswd", "^\\.?htpasswd$", severity.Low, t)
	shouldFail("htpasswd", "^\\.?htpasswd$", severity.Low, t)
	shouldFail(".netrc", "^\\.?netrc$", severity.Low, t)
	shouldFail("netrc", "^\\.?netrc$", severity.Low, t)
	shouldFail("foo.tblk", "^.*\\.tblk$", severity.Low, t) //Tunnelblick
	shouldFail("foo.ovpn", "^.*\\.ovpn$", severity.Low, t) //OpenVPN
}

func TestShouldFlagPotentialCrendentialDatabases(t *testing.T) {
	shouldFail("foo.kdb", "^.*\\.kdb$", severity.Low, t)                     //KeePass
	shouldFail("foo.agilekeychain", "^.*\\.agilekeychain$", severity.Low, t) //1Password
	shouldFail("foo.keychain", "^.*\\.keychain$", severity.Low, t)           //apple keychain
	shouldFail("foo.keystore", "^.*\\.key(store|ring)$", severity.Low, t)    //gnome keyring db
	shouldFail("foo.keyring", "^.*\\.key(store|ring)$", severity.Low, t)     //gnome keyring db
}

func TestShouldFlagPotentialJenkinsAndCICompromises(t *testing.T) {
	shouldFail("jenkins.plugins.publish_over_ssh.BapSshPublisherPlugin.xml", "^jenkins\\.plugins\\.publish_over_ssh\\.BapSshPublisherPlugin.xml$", severity.Low, t)
	shouldFail("credentials.xml", "^credentials\\.xml$", severity.Low, t)
	shouldFail("foo.pubxml.user", "^.*\\.pubxml(\\.user)?$", severity.Low, t)
	shouldFail("foo.pubxml", "^.*\\.pubxml(\\.user)?$", severity.Low, t)
}

func TestShouldFlagPotentialConfigurationFilesThatMightContainSensitiveInformation(t *testing.T) {
	shouldFail(".s3cfg", "^\\.?s3cfg$", severity.Low, t)      //s3 configuration
	shouldFail("foo.ovpn", "^.*\\.ovpn$", severity.Low, t)    //OpenVPN configuration
	shouldFail(".gitrobrc", "^\\.gitrobrc$", severity.Low, t) //Gitrob configuration
	shouldFail(".bashrc", "^\\.?(bash|zsh)rc$", severity.Low, t)
	shouldFail(".zshrc", "^\\.?(bash|zsh)rc$", severity.Low, t)
	shouldFail(".profile", "^\\.?(bash_|zsh_)?profile$", severity.Low, t)
	shouldFail(".bash_profile", "^\\.?(bash_|zsh_)?profile$", severity.Low, t)
	shouldFail(".zsh_profile", "^\\.?(bash_|zsh_)?profile$", severity.Low, t)
	shouldFail(".bash_aliases", "^\\.?(bash_|zsh_)?aliases$", severity.Low, t)
	shouldFail(".zsh_aliases", "^\\.?(bash_|zsh_)?aliases$", severity.Low, t)
	shouldFail(".aliases", "^\\.?(bash_|zsh_)?aliases$", severity.Low, t)
	shouldFail("secret_token.rb", "^secret_token.rb$", severity.Low, t)          //Rails secret token. http://www.exploit-db.com/exploits/27527
	shouldFail("omniauth.rb", "^omniauth.rb$", severity.Low, t)                  //OmniAuth configuration file, client application secrets
	shouldFail("carrierwave.rb", "^carrierwave.rb$", severity.Low, t)            //May contain Amazon S3 and Google Storage credentials
	shouldFail("schema.rb", "^schema.rb$", severity.Low, t)                      //Rails application DB schema info
	shouldFail("database.yml", "^database.yml$", severity.Low, t)                //Rails db connection strings
	shouldFail("settings.py", "^settings.py$", severity.Low, t)                  //Django credentials, keys etc
	shouldFail("wp-config.php", "^.*(config)(\\.inc)?\\.php$", severity.Low, t)  //Wordpress PHP config file
	shouldFail("config.php", "^.*(config)(\\.inc)?\\.php$", severity.Low, t)     //General PHP config file
	shouldFail("config.inc.php", "^.*(config)(\\.inc)?\\.php$", severity.Low, t) //PHP MyAdmin file with credentials etc
	shouldFail("LocalSettings.php", "^LocalSettings.php$", severity.Low, t)      //MediaWiki configuration file
	shouldFail(".env", "\\.?env", severity.Low, t)                               //PHP environment file that contains sensitive data
}

func TestShouldFlagPotentialSuspiciousSoundingFileNames(t *testing.T) {
	shouldFail("database.dump", "\\bdump|dump\\b", severity.Low, t) //Dump might contain sensitive information
	shouldFail("foo.sql", "\\bsql|sql\\b", severity.Low, t)         //Sql file, might be a dump and contain sensitive information
	shouldFail("mydb.sqldump", "\\bdump|dump\\b", severity.Low, t)  //Sql file, dump file, might be a dump and contain sensitive information

	shouldFail("foo_password", "password", severity.Low, t)     //Looks like a password?
	shouldFail("foo.password", "password", severity.Low, t)     //Looks like a password?
	shouldFail("foo_password.txt", "password", severity.Low, t) //Looks like a password?

	shouldFail("foo_backup", "backup", severity.Low, t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo.backup", "backup", severity.Low, t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo_backup.txt", "backup", severity.Low, t) //Looks like a backup. Might contain sensitive information.

	shouldFail("private_key", "private.*key", severity.Low, t)     //Looks like a private key.
	shouldFail("private.key", "private.*key", severity.Low, t)     //Looks like a private key.
	shouldFail("private_key.txt", "private.*key", severity.Low, t) //Looks like a private key.
	shouldFail("otr.private_key", "private.*key", severity.Low, t)

	shouldFail("oauth_token", "(oauth).*(token)", severity.Low, t)     //Looks like an oauth token
	shouldFail("oauth.token", "(oauth).*(token)", severity.Low, t)     //Looks like an oauth token
	shouldFail("oauth_token.txt", "(oauth).*(token)", severity.Low, t) //Looks like an oauth token

	shouldFail("development.log", "^.*\\.log$", severity.Low, t) //Looks like a log file, could contain sensitive information
}

func TestFilenameDetectorReportsFailuresIfAnyFileInAdditionsMatchesAnyFlagPattern(t *testing.T) {
	shouldFail(".kwallet", "^\\.?kwallet$", severity.Low, t)
	shouldFail("kwallet", "^\\.?kwallet$", severity.Low, t)
	shouldFail(".gnucash", "^\\.?gnucash$", severity.Low, t)
	shouldFail("gnucash", "^\\.?gnucash$", severity.Low, t)
}

func TestShouldIgnoreFilesWhenAskedToDoSoByIgnores(t *testing.T) {
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "id_rsa", severity.Low, t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "*_rsa", severity.Low, t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_dsa", "id_*", severity.Low, t)
}

func TestShouldIgnoreIfErrorIsBelowThreshold(t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	severity := severity.High
	fileName := ".bash_aliases"
	DefaultFileNameDetector(severity).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results, func() {})
	assert.False(t, results.HasFailures(), "Expected file %s to not fail", fileName)
	assert.True(t, results.HasWarnings(), "Expected file %s to having warnings", fileName)
}

func shouldFail(fileName, pattern string, threshold severity.Severity, t *testing.T) {
	shouldFailWithSpecificPattern(fileName, pattern, threshold, t)
	shouldFailWithDefaultDetector(fileName, pattern, threshold, t)
}

func shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors(fileName, ignore string, threshold severity.Severity, t *testing.T) {
	shouldFailWithDefaultDetector(fileName, "", threshold, t)
	shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore, threshold, t)
}

func shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore string, threshold severity.Severity, t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)

	fileIgnoreConfig := talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = ignore
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filename"
	talismanRC := &talismanrc.TalismanRC{}
	talismanRC.FileIgnoreConfig = make([]talismanrc.FileIgnoreConfig, 1)
	talismanRC.FileIgnoreConfig[0] = fileIgnoreConfig

	DefaultFileNameDetector(threshold).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results, func() {})
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", fileName, ignore)
}

func shouldFailWithSpecificPattern(fileName, pattern string, threshold severity.Severity, t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	pt := []*severity.PatternSeverity{{Pattern: regexp.MustCompile(pattern), Severity: severity.Low}}
	NewFileNameDetector(pt, threshold).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file %s to fail the check against the %s pattern", fileName, pattern)
}

func shouldFailWithDefaultDetector(fileName, pattern string, severity severity.Severity, t *testing.T) {
	results := helpers.NewDetectionResults(talismanrc.Hook)
	DefaultFileNameDetector(severity).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results, func() {})
	assert.True(t, results.HasFailures(), "Expected file %s to fail the check against default detector. Missing pattern %s?", fileName, pattern)
}

func additionsNamed(names ...string) []gitrepo.Addition {
	result := make([]gitrepo.Addition, len(names))
	for i, name := range names {
		result[i] = gitrepo.Addition{
			Path: gitrepo.FilePath(name),
			Name: gitrepo.FileName(name),
			Data: make([]byte, 0),
		}
	}
	return result
}
