package detector

import (
	"testing"

	"talisman/gitrepo"

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnoreEmptyLinesInTheFile(t *testing.T) {
	for _, s := range []string{"", " ", "  "} {
		assert.True(t, NewTalismanRCIgnore([]byte(s)).AcceptsAll(), "Expected '%s' to result in no ignore patterns.", s)
	}
}

func TestShouldIgnoreUnformattedFiles(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		assert.True(t, NewTalismanRCIgnore([]byte(s)).AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
	}
}

func TestShouldParseIgnoreLinesProperly(t *testing.T) {
	assert.Equal(t, NewIgnores("foo* # comment"), SingleIgnore("foo*", "comment"))
	assert.Equal(t, NewIgnores("foo* # comment with multiple words"), SingleIgnore("foo*", "comment with multiple words"))
	assert.Equal(t, NewIgnores("foo* # comment with#multiple#words"), SingleIgnore("foo*", "comment with#multiple#words"))
	assert.Equal(t, NewIgnores("foo*# comment"), SingleIgnore("foo*", "comment"))

	assert.Equal(t, NewIgnores("# comment"), SingleIgnore("", "comment"))
	assert.Equal(t, NewIgnores("#comment"), SingleIgnore("", "comment"))

	assert.Equal(t, NewIgnores(""), SingleIgnore("", ""))
	assert.Equal(t, NewIgnores(" "), SingleIgnore("", ""))

	assert.Equal(t, NewIgnores("foo # ignore:some-detector"), SingleIgnore("foo", "ignore:some-detector", "some-detector"))
	assert.Equal(t, NewIgnores("foo # ignore:some-detector,some-other-detector"), SingleIgnore("foo", "ignore:some-detector,some-other-detector", "some-detector", "some-other-detector"))
	assert.Equal(t, NewIgnores("foo # ignore:some-detector because of some reason"), SingleIgnore("foo", "ignore:some-detector because of some reason", "some-detector"))

}

//Need to work on this test case as it deals with matching file names throughout directories and sub directories which talismanrc does not yet support
// func TestRawPatterns(t *testing.T) {
// 	assertAccepts("foo", "", "bar", t)
// 	assertAccepts("foo", "", "foobar", t)
// 	assertAccepts("foo", "", "foo/bar", t)
// 	assertAccepts("foo", "", "foo/bar/baz", t)

// 	assertDenies("foo", "", "foo", t)
// 	assertDenies("foo", "", "bar/foo", t)
// 	assertDenies("foo", "", "bar/baz/foo", t)
// }

//Need to work on this test case as it deals with wildcards and talismanrc does not yet support wildcards
// func TestSingleStarPatterns(t *testing.T) {
// 	assertAccepts("foo*", "", "bar", t)
// 	assertAccepts("foo*", "", "foo/bar", t)
// 	assertAccepts("foo*", "", "foo/bar/baz", t)

// 	assertDenies("foo*", "", "foo", t)
// 	assertDenies("foo*", "", "foobar", t)

// 	assertAccepts("*.pem", "", "foo.txt", t)
// 	assertAccepts("*.pem", "", "pem.go", t)
// 	assertDenies("*.pem", "", "secret.pem", t)
// 	assertDenies("*.pem", "", "foo/bar/secret.pem", t)
// }

func TestDirectoryPatterns(t *testing.T) {
	assertAccepts("foo/", "", "bar", t)
	assertAccepts("foo/", "", "foo", t)
	assertDenies("foo/", "filename", "foo/bar", t)
	assertDenies("foo/", "filename", "foo/bar.txt", t)
	assertDenies("foo/", "filename", "foo/bar/baz.txt", t)
}

func TestIgnoreAdditionsByScope(t *testing.T) {
	file1 := testAddition("yarn.lock")
	file2 := testAddition("similaryarn.lock")
	file3 := testAddition("java.lock")
	file4 := testAddition("Gopkg.lock")
	file5 := testAddition("vendors/abc")
	file6 := testAddition(".idea/asd")
	additions := []git_repo.Addition{file1, file2, file3, file4, file5, file6}

	scopesToIgnore := []string{"node", "go", "idea"}
	talismanRCIgnoreConfig := CreateTalismanRCIgnoreWithScopeIgnore(scopesToIgnore)

	nodeIgnores := []string{"node.lock", "*yarn.lock"}
	javaIgnores := []string{"java.lock"}
	goIgnores := []string{"go.lock", "Gopkg.lock", "vendors/"}
	ideaIgnores := []string{".idea/"}
	scopesMap := map[string][]string{"node": nodeIgnores, "java": javaIgnores, "go": goIgnores, "idea": ideaIgnores}

	filteredAdditions := IgnoreAdditionsByScope(additions, talismanRCIgnoreConfig, scopesMap)

	assert.NotContains(t, filteredAdditions, file1)
	assert.NotContains(t, filteredAdditions, file2)
	assert.Contains(t, filteredAdditions, file3)
	assert.NotContains(t, filteredAdditions, file4)
	assert.NotContains(t, filteredAdditions, file5)
	assert.NotContains(t, filteredAdditions, file6)
}

//Need to work on this test case as it deals with comments and talismanrc does not deal in comments
//func TestCommentPatterns(t *testing.T) {
//	assertAccepts("foo # some comment", "bar", t)
//	assertDenies("foo # some comment", "foo", t)
//
//	assertDenies("foo* # comment", "foo", t)
//	assertAccepts("foo* # comment", "bar", t)
//
//	assertAccepts("foo/ # comment", "bar", t)
//	assertDenies("foo/ # comment", "foo/bar", t)
//}

func TestIgnoringDetectors(t *testing.T) {
	assertDeniesDetector("foo", "someDetector", "foo", "someDetector", t)
	assertAcceptsDetector("foo", "someDetector", "foo", "someOtherDetector", t)
}

func assertDenies(line, ignoreDetector string, path string, t *testing.T) {
	assertDeniesDetector(line, ignoreDetector, path, "filename", t)
}

func assertDeniesDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, CreateTalismanRCIgnoreWithFileName(line, ignoreDetector).Deny(testAddition(path), detectorName), "%s is expected to deny a file named %s.", line, path)
}

func assertAccepts(line, ignoreDetector string, path string, t *testing.T, detectorNames ...string) {
	assertAcceptsDetector(line, ignoreDetector, path, "someDetector", t)
}

func assertAcceptsDetector(line, ignoreDetector string, path string, detectorName string, t *testing.T) {
	assert.True(t, CreateTalismanRCIgnoreWithFileName(line, ignoreDetector).Accept(testAddition(path), detectorName), "%s is expected to accept a file named %s.", line, path)
}

func testAddition(path string) gitrepo.Addition {
	return gitrepo.NewAddition(path, make([]byte, 0))
}

func CreateTalismanRCIgnoreWithFileName(filename string, detector string) TalismanRCIgnore {
	fileIgnoreConfig := FileIgnoreConfig{}
	fileIgnoreConfig.FileName = filename
	if detector != "" {
		fileIgnoreConfig.IgnoreDetectors = make([]string, 1)
		fileIgnoreConfig.IgnoreDetectors[0] = detector
	}
	talismanRCIgnore := TalismanRCIgnore{}
	talismanRCIgnore.FileIgnoreConfig = make([]FileIgnoreConfig, 1)
	talismanRCIgnore.FileIgnoreConfig[0] = fileIgnoreConfig
	return talismanRCIgnore
}

func CreateTalismanRCIgnoreWithScopeIgnore(scopesToIgnore []string) TalismanRCIgnore {
	var scopeConfigs []ScopeConfig
	for _, scopeIgnore := range scopesToIgnore {
		scopeIgnoreConfig := ScopeConfig{}
		scopeIgnoreConfig.ScopeName = scopeIgnore
		scopeConfigs = append(scopeConfigs, scopeIgnoreConfig)
	}

	talismanRCIgnore := TalismanRCIgnore{ScopeConfig: scopeConfigs}
	return talismanRCIgnore
}

func SingleIgnore(pattern string, comment string, ignoredDetectors ...string) Ignores {
	return Ignores{patterns: []Ignore{{
		pattern:          pattern,
		comment:          comment,
		ignoredDetectors: ignoredDetectors,
	}}}
}
