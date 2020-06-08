package detector

//This is completely derived from the really useful work done by Jen Andre here:
//https://github.com/jandre/safe-commit-hook

import (
	"regexp"
	"testing"

	"talisman/gitrepo"
	"talisman/talismanrc"

	"github.com/stretchr/testify/assert"
)

func TestShouldFlagPotentialSSHPrivateKeys(t *testing.T) {
	shouldFail("id_rsa", "^.+_rsa$", t)
	shouldFail("id_dsa", "^.+_dsa.*$", t)
	shouldFail("id_dsa.pub", "^.+_dsa.*$", t)
	shouldFail("id_ed25519", "^.+_ed25519$", t)
	shouldFail("id_ecdsa", "^.+_ecdsa$", t)
}

func TestShouldFlagPotentialHistoryFiles(t *testing.T) {
	shouldFail(".bash_history", "^\\.\\w+_history$", t)
	shouldFail(".zsh_history", "^\\.\\w+_history$", t)
	shouldFail(".z_history", "^\\.\\w+_history$", t)
	shouldFail(".irb_history", "^\\.\\w+_history$", t)
	shouldFail(".psql_history", "^\\.\\w+_history$", t)
	shouldFail(".mysql_history", "^\\.\\w+_history$", t)
}

func TestShouldFlagPotentialPrivateKeys(t *testing.T) {
	shouldFail("foo.pem", "^.+\\.pem$", t)
	shouldFail("foo.ppk", "^.+\\.ppk$", t)
	shouldFail("foo.key", "^.+\\.key(pair)?$", t)
	shouldFail("foo.keypair", "^.+\\.key(pair)?$", t)
}

func TestShouldFlagPotentialKeyBundles(t *testing.T) {
	shouldFail("foo.pkcs12", "^.+\\.pkcs12$", t)
	shouldFail("foo.pfx", "^.+\\.pfx$", t)
	shouldFail("foo.p12", "^.+\\.p12$", t)
	shouldFail("foo.asc", "^.+\\.asc$", t)
}

func TestShouldFlagPotentialConfigurationFiles(t *testing.T) {
	shouldFail(".htpasswd", "^\\.?htpasswd$", t)
	shouldFail("htpasswd", "^\\.?htpasswd$", t)
	shouldFail(".netrc", "^\\.?netrc$", t)
	shouldFail("netrc", "^\\.?netrc$", t)
	shouldFail("foo.tblk", "^.*\\.tblk$", t) //Tunnelblick
	shouldFail("foo.ovpn", "^.*\\.ovpn$", t) //OpenVPN
}

func TestShouldFlagPotentialCrendentialDatabases(t *testing.T) {
	shouldFail("foo.kdb", "^.*\\.kdb$", t)                     //KeePass
	shouldFail("foo.agilekeychain", "^.*\\.agilekeychain$", t) //1Password
	shouldFail("foo.keychain", "^.*\\.keychain$", t)           //apple keychain
	shouldFail("foo.keystore", "^.*\\.key(store|ring)$", t)    //gnome keyring db
	shouldFail("foo.keyring", "^.*\\.key(store|ring)$", t)     //gnome keyring db
}

func TestShouldFlagPotentialJenkinsAndCICompromises(t *testing.T) {
	shouldFail("jenkins.plugins.publish_over_ssh.BapSshPublisherPlugin.xml", "^jenkins\\.plugins\\.publish_over_ssh\\.BapSshPublisherPlugin.xml$", t)
	shouldFail("credentials.xml", "^credentials\\.xml$", t)
	shouldFail("foo.pubxml.user", "^.*\\.pubxml(\\.user)?$", t)
	shouldFail("foo.pubxml", "^.*\\.pubxml(\\.user)?$", t)
}

func TestShouldFlagPotentialConfigurationFilesThatMightContainSensitiveInformation(t *testing.T) {
	shouldFail(".s3cfg", "^\\.?s3cfg$", t)      //s3 configuration
	shouldFail("foo.ovpn", "^.*\\.ovpn$", t)    //OpenVPN configuration
	shouldFail(".gitrobrc", "^\\.gitrobrc$", t) //Gitrob configuration
	shouldFail(".bashrc", "^\\.?(bash|zsh)rc$", t)
	shouldFail(".zshrc", "^\\.?(bash|zsh)rc$", t)
	shouldFail(".profile", "^\\.?(bash_|zsh_)?profile$", t)
	shouldFail(".bash_profile", "^\\.?(bash_|zsh_)?profile$", t)
	shouldFail(".zsh_profile", "^\\.?(bash_|zsh_)?profile$", t)
	shouldFail(".bash_aliases", "^\\.?(bash_|zsh_)?aliases$", t)
	shouldFail(".zsh_aliases", "^\\.?(bash_|zsh_)?aliases$", t)
	shouldFail(".aliases", "^\\.?(bash_|zsh_)?aliases$", t)
	shouldFail("secret_token.rb", "^secret_token.rb$", t)          //Rails secret token. http://www.exploit-db.com/exploits/27527
	shouldFail("omniauth.rb", "^omniauth.rb$", t)                  //OmniAuth configuration file, client application secrets
	shouldFail("carrierwave.rb", "^carrierwave.rb$", t)            //May contain Amazon S3 and Google Storage credentials
	shouldFail("schema.rb", "^schema.rb$", t)                      //Rails application DB schema info
	shouldFail("database.yml", "^database.yml$", t)                //Rails db connection strings
	shouldFail("settings.py", "^settings.py$", t)                  //Django credentials, keys etc
	shouldFail("wp-config.php", "^.*(config)(\\.inc)?\\.php$", t)  //Wordpress PHP config file
	shouldFail("config.php", "^.*(config)(\\.inc)?\\.php$", t)     //General PHP config file
	shouldFail("config.inc.php", "^.*(config)(\\.inc)?\\.php$", t) //PHP MyAdmin file with credentials etc
	shouldFail("LocalSettings.php", "^LocalSettings.php$", t)      //MediaWiki configuration file
	shouldFail(".env", "\\.?env", t)                               //PHP environment file that contains sensitive data
}

func TestShouldFlagPotentialSuspiciousSoundingFileNames(t *testing.T) {
	shouldFail("database.dump", "\\bdump|dump\\b", t) //Dump might contain sensitive information
	shouldFail("foo.sql", "\\bsql|sql\\b", t)         //Sql file, might be a dump and contain sensitive information
	shouldFail("mydb.sqldump", "\\bdump|dump\\b", t)  //Sql file, dump file, might be a dump and contain sensitive information

	shouldFail("foo_password", "password", t)     //Looks like a password?
	shouldFail("foo.password", "password", t)     //Looks like a password?
	shouldFail("foo_password.txt", "password", t) //Looks like a password?

	shouldFail("foo_backup", "backup", t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo.backup", "backup", t)     //Looks like a backup. Might contain sensitive information.
	shouldFail("foo_backup.txt", "backup", t) //Looks like a backup. Might contain sensitive information.

	shouldFail("private_key", "private.*key", t)     //Looks like a private key.
	shouldFail("private.key", "private.*key", t)     //Looks like a private key.
	shouldFail("private_key.txt", "private.*key", t) //Looks like a private key.
	shouldFail("otr.private_key", "private.*key", t)

	shouldFail("oauth_token", "(oauth).*(token)", t)     //Looks like an oauth token
	shouldFail("oauth.token", "(oauth).*(token)", t)     //Looks like an oauth token
	shouldFail("oauth_token.txt", "(oauth).*(token)", t) //Looks like an oauth token

	shouldFail("development.log", "^.*\\.log$", t) //Looks like a log file, could contain sensitive information
}

func TestFilenameDetectorReportsFailuresIfAnyFileInAdditionsMatchesAnyFlagPattern(t *testing.T) {
	shouldFail(".kwallet", "^\\.?kwallet$", t)
	shouldFail("kwallet", "^\\.?kwallet$", t)
	shouldFail(".gnucash", "^\\.?gnucash$", t)
	shouldFail("gnucash", "^\\.?gnucash$", t)
}

func TestShouldIgnoreFilesWhenAskedToDoSoByIgnores(t *testing.T) {
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "id_rsa", t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_rsa", "*_rsa", t)
	shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors("id_dsa", "id_*", t)
}

func shouldFail(fileName, pattern string, t *testing.T) {
	shouldFailWithSpecificPattern(fileName, pattern, t)
	shouldFailWithDefaultDetector(fileName, pattern, t)
}

func shouldIgnoreFilesWhichWouldOtherwiseTriggerErrors(fileName, ignore string, t *testing.T) {
	shouldFailWithDefaultDetector(fileName, "", t)
	shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore, t)
}

func shouldNotFailWithDefaultDetectorAndIgnores(fileName, ignore string, t *testing.T) {
	results := NewDetectionResults()

	fileIgnoreConfig := talismanrc.FileIgnoreConfig{}
	fileIgnoreConfig.FileName = ignore
	fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
	fileIgnoreConfig.IgnoreDetectors[0] = "filename"
	talismanRC := &talismanrc.TalismanRC{}
	talismanRC.FileIgnoreConfig = make([]talismanrc.FileIgnoreConfig, 1)
	talismanRC.FileIgnoreConfig[0] = fileIgnoreConfig

	DefaultFileNameDetector().Test(ChecksumCompare{calculator: nil, talismanRC: talismanrc.NewTalismanRC(nil)}, additionsNamed(fileName), talismanRC, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", fileName, ignore)
}

func shouldFailWithSpecificPattern(fileName, pattern string, t *testing.T) {
	results := NewDetectionResults()
	pt := regexp.MustCompile(pattern)
	NewFileNameDetector([]*regexp.Regexp{pt}).Test(ChecksumCompare{calculator: nil, talismanRC: talismanrc.NewTalismanRC(nil)}, additionsNamed(fileName), talismanRC, results)
	assert.True(t, results.HasFailures(), "Expected file %s to fail the check against the %s pattern", fileName, pattern)
}

func shouldFailWithDefaultDetector(fileName, pattern string, t *testing.T) {
	results := NewDetectionResults()
	DefaultFileNameDetector().Test(ChecksumCompare{calculator: nil, talismanRC: talismanrc.NewTalismanRC(nil)}, additionsNamed(fileName), talismanRC, results)
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
