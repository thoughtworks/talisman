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
	shouldFail("id_rsa", "^.+_rsa$", severity.LowSeverity, t)
	shouldFail("id_dsa", "^.+_dsa.*$", severity.LowSeverity, t)
	shouldFail("id_dsa.pub", "^.+_dsa.*$", severity.LowSeverity, t)
	shouldFail("id_ed25519", "^.+_ed25519$", severity.LowSeverity, t)
	shouldFail("id_ecdsa", "^.+_ecdsa$", severity.LowSeverity, t)
}

func TestShouldFlagPotentialHistoryFiles(t *testing.T) {
	shouldFail(".bash_history", "^\\.\\w+_history$", severity.LowSeverity, t)
	shouldFail(".zsh_history", "^\\.\\w+_history$", severity.LowSeverity, t)
	shouldFail(".z_history", "^\\.\\w+_history$", severity.LowSeverity, t)
	shouldFail(".irb_history", "^\\.\\w+_history$", severity.LowSeverity, t)
	shouldFail(".psql_history", "^\\.\\w+_history$", severity.LowSeverity, t)
	shouldFail(".mysql_history", "^\\.\\w+_history$", severity.LowSeverity, t)
}

func TestShouldFlagPotentialPrivateKeys(t *testing.T) {
	shouldFail("foo.pem", "^.+\\.pem$", severity.LowSeverity, t)
	shouldFail("foo.ppk", "^.+\\.ppk$", severity.LowSeverity, t)
	shouldFail("foo.key", "^.+\\.key(pair)?$", severity.LowSeverity, t)
	shouldFail("foo.keypair", "^.+\\.key(pair)?$", severity.LowSeverity, t)
}

func TestShouldFlagPotentialKeyBundles(t *testing.T) {
	shouldFail("foo.pkcs12", "^.+\\.pkcs12$", severity.LowSeverity, t)
	shouldFail("foo.pfx", "^.+\\.pfx$", severity.LowSeverity, t)
	shouldFail("foo.p12", "^.+\\.p12$", severity.LowSeverity, t)
	shouldFail("foo.asc", "^.+\\.asc$", severity.LowSeverity, t)
}

func TestShouldFlagPotentialConfigurationFiles(t *testing.T) {
	shouldFail(".htpasswd", "^\\.?htpasswd$", severity.LowSeverity, t)
	shouldFail("htpasswd", "^\\.?htpasswd$", severity.LowSeverity, t)
	shouldFail(".netrc", "^\\.?netrc$", severity.LowSeverity, t)
	shouldFail("netrc", "^\\.?netrc$", severity.LowSeverity, t)
	shouldFail("foo.tblk", "^.*\\.tblk$", severity.LowSeverity, t) //Tunnelblick
	shouldFail("foo.ovpn", "^.*\\.ovpn$", severity.LowSeverity, t) //OpenVPN
}

func TestShouldFlagPotentialCrendentialDatabases(t *testing.T) {
	shouldFail("foo.kdb", "^.*\\.kdb$", severity.LowSeverity, t)                     //KeePass
	shouldFail("foo.agilekeychain", "^.*\\.agilekeychain$", severity.LowSeverity, t) //1Password
	shouldFail("foo.keychain", "^.*\\.keychain$", severity.LowSeverity, t)           //apple keychain
	shouldFail("foo.keystore", "^.*\\.key(store|ring)$", severity.LowSeverity, t)    //gnome keyring db
	shouldFail("foo.keyring", "^.*\\.key(store|ring)$", severity.LowSeverity, t)     //gnome keyring db
}

func TestShouldFlagPotentialJenkinsAndCICompromises(t *testing.T) {
	shouldFail("jenkins.plugins.publish_over_ssh.BapSshPublisherPlugin.xml", "^jenkins\\.plugins\\.publish_over_ssh\\.BapSshPublisherPlugin.xml$", severity.LowSeverity, t)
	shouldFail("credentials.xml", "^credentials\\.xml$", severity.LowSeverity, t)
	shouldFail("foo.pubxml.user", "^.*\\.pubxml(\\.user)?$", severity.LowSeverity, t)
	shouldFail("foo.pubxml", "^.*\\.pubxml(\\.user)?$", severity.LowSeverity, t)
}

func TestShouldFlagPotentialConfigurationFilesThatMightContainSensitiveInformation(t *testing.T) {
	shouldFail(".s3cfg", "^\\.?s3cfg$", severity.LowSeverity, t)      //s3 configuration
	shouldFail("foo.ovpn", "^.*\\.ovpn$", severity.LowSeverity, t)    //OpenVPN configuration
	shouldFail(".gitrobrc", "^\\.gitrobrc$", severity.LowSeverity, t) //Gitrob configuration
	shouldFail(".bashrc", "^\\.?(bash|zsh)rc$", severity.LowSeverity, t)
	shouldFail(".zshrc", "^\\.?(bash|zsh)rc$", severity.LowSeverity, t)
	shouldFail(".profile", "^\\.?(bash_|zsh_)?profile$", severity.LowSeverity, t)
	shouldFail(".bash_profile", "^\\.?(bash_|zsh_)?profile$", severity.LowSeverity, t)
	shouldFail(".zsh_profile", "^\\.?(bash_|zsh_)?profile$", severity.LowSeverity, t)
	shouldFail(".bash_aliases", "^\\.?(bash_|zsh_)?aliases$", severity.LowSeverity, t)
	shouldFail(".zsh_aliases", "^\\.?(bash_|zsh_)?aliases$", severity.LowSeverity, t)
	shouldFail(".aliases", "^\\.?(bash_|zsh_)?aliases$", severity.LowSeverity, t)
	shouldFail("secret_token.rb", "^secret_token.rb$", severity.LowSeverity, t)          //Rails secret token. http://www.exploit-db.com/exploits/27527
	shouldFail("omniauth.rb", "^omniauth.rb$", severity.LowSeverity, t)                  //OmniAuth configuration file, client application secrets
	shouldFail("carrierwave.rb", "^carrierwave.rb$", severity.LowSeverity, t)            //May contain Amazon S3 and Google Storage credentials
	shouldFail("schema.rb", "^schema.rb$", severity.LowSeverity, t)                      //Rails application DB schema info
	shouldFail("database.yml", "^database.yml$", severity.LowSeverity, t)                //Rails db connection strings
	shouldFail("settings.py", "^settings.py$", severity.LowSeverity, t)                  //Django credentials, keys etc
	shouldFail("wp-config.php", "^.*(config)(\\.inc)?\\.php$", severity.LowSeverity, t)  //Wordpress PHP config file
	shouldFail("config.php", "^.*(config)(\\.inc)?\\.php$", severity.LowSeverity, t)     //General PHP config file
	shouldFail("config.inc.php", "^.*(config)(\\.inc)?\\.php$", severity.LowSeverity, t) //PHP MyAdmin file with credentials etc
	shouldFail("LocalSettings.php", "^LocalSettings.php$", severity.LowSeverity, t)      //MediaWiki configuration file
	shouldFail(".env", "\\.?env", severity.LowSeverity, t)                               //PHP environment file that contains sensitive data
}

func TestShouldFlagPotentialSuspiciousSoundingFileNames(t *testing.T) {
	shouldFail("database.dump", "\\bdump|dump\\b", severity.LowSeverity, t) //Dump might contain sensitive information
	shouldFail("foo.sql", "\\bsql|sql\\b", severity.LowSeverity, t)         //Sql file, might be a dump and contain sensitive information
	shouldFail("mydb.sqldump", "\\bdump|dump\\b", severity.LowSeverity, t)  //Sql file, dump file, might be a dump and contain sensitive information

	shouldFail("foo_password", "password", severity.LowSeverity, t)     //Looks like a password?
	shouldFail("foo.password", "password", severity.LowSeverity, t)     //Looks like a password?
	shouldFail("foo_password.txt", "password", severity.LowSeverity, t) //Looks like a password?

	shouldFail("foo_backup", "backup", severity.LowSeverity, t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo.backup", "backup", severity.LowSeverity, t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo_backup.txt", "backup", severity.LowSeverity, t) //Looks like a backup. Might contain sensitive information.

	shouldFail("private_key", "private.*key", severity.LowSeverity, t)     //Looks like a private key.
	shouldFail("private.key", "private.*key", severity.LowSeverity, t)     //Looks like a private key.
	shouldFail("private_key.txt", "private.*key", severity.LowSeverity, t) //Looks like a private key.
	shouldFail("otr.private_key", "private.*key", severity.LowSeverity, t)

	shouldFail("oauth_token", "(oauth).*(token)", severity.LowSeverity, t)     //Looks like an oauth token
	shouldFail("oauth.token", "(oauth).*(token)", severity.LowSeverity, t)     //Looks like an oauth token
	shouldFail("oauth_token.txt", "(oauth).*(token)", severity.LowSeverity, t) //Looks like an oauth token

	shouldFail("development.log", "^.*\\.log$", severity.LowSeverity, t) //Looks like a log file, could contain sensitive information
}

func TestFilenameDetectorReportsFailuresIfAnyFileInAdditionsMatchesAnyFlagPattern(t *testing.T) {
	shouldFail(".kwallet", "^\\.?kwallet$", severity.LowSeverity, t)
	shouldFail("kwallet", "^\\.?kwallet$", severity.LowSeverity, t)
	shouldFail(".gnucash", "^\\.?gnucash$", severity.LowSeverity, t)
	shouldFail("gnucash", "^\\.?gnucash$", severity.LowSeverity, t)
}

func TestShouldIgnoreFilesWhenAskedToDoSoByIgnores(t *testing.T) {
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "id_rsa", severity.LowSeverity, t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "*_rsa", severity.LowSeverity, t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_dsa", "id_*", severity.LowSeverity, t)
}

func TestShouldIgnoreIfErrorIsBelowThreshold(t *testing.T) {
	results := helpers.NewDetectionResults()
	severity := severity.HighSeverity
	fileName := ".bash_aliases"
	DefaultFileNameDetector(severity).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results)
	assert.False(t, results.HasFailures(), "Expected file %s to not fail", fileName)
	assert.True(t, results.HasWarnings(), "Expected file %s to having warnings", fileName)
}

func shouldFail(fileName, pattern string, threshold severity.SeverityValue, t *testing.T) {
	shouldFailWithSpecificPattern(fileName, pattern, threshold, t)
	shouldFailWithDefaultDetector(fileName, pattern, threshold, t)
}

func shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors(fileName, ignore string, threshold severity.SeverityValue, t *testing.T) {
	shouldFailWithDefaultDetector(fileName, "", threshold, t)
	shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore, threshold, t)
}

func shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore string, threshold severity.SeverityValue, t *testing.T) {
	results := helpers.NewDetectionResults()

	fileIgnoreConfig := talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = ignore
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filename"
	talismanRC := &talismanrc.TalismanRC{}
	talismanRC.FileIgnoreConfig = make([]talismanrc.FileIgnoreConfig, 1)
	talismanRC.FileIgnoreConfig[0] = fileIgnoreConfig

	DefaultFileNameDetector(threshold).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", fileName, ignore)
}

func shouldFailWithSpecificPattern(fileName, pattern string, threshold severity.SeverityValue, t *testing.T) {
	results := helpers.NewDetectionResults()
	pt := []*severity.PatternSeverity{{Pattern: regexp.MustCompile(pattern), Severity: severity.Low()}}
	NewFileNameDetector(pt, threshold).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results)
	assert.True(t, results.HasFailures(), "Expected file %s to fail the check against the %s pattern", fileName, pattern)
}

func shouldFailWithDefaultDetector(fileName, pattern string, severity severity.SeverityValue, t *testing.T) {
	results := helpers.NewDetectionResults()
	DefaultFileNameDetector(severity).Test(helpers.NewChecksumCompare(nil, utility.DefaultSHA256Hasher{}, talismanrc.NewTalismanRC(nil)), additionsNamed(fileName), talismanRC, results)
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
