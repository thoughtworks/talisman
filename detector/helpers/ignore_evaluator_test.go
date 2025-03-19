package helpers

import (
	"io"
	"talisman/gitrepo"
	mockchecksumcalculator "talisman/internal/mock/checksumcalculator"
	"talisman/talismanrc"

	"github.com/golang/mock/gomock"
	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
)

func init() {
	logr.SetOutput(io.Discard)
}
func TestIsScanNotRequired(t *testing.T) {

	t.Run("should return false if talismanrc is empty", func(t *testing.T) {
		ignoreConfig := &talismanrc.TalismanRC{
			IgnoreConfigs: []talismanrc.IgnoreConfig{},
		}
		ie := IgnoreEvaluator{nil, ignoreConfig}
		addition := gitrepo.Addition{Path: "some.txt"}

		required := ie.isScanNotRequired(addition)

		assert.False(t, required)
	})

	t.Run("should loop through talismanrc configs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		checksumCalculator := mockchecksumcalculator.NewMockChecksumCalculator(ctrl)
		ignoreConfig := talismanrc.TalismanRC{
			IgnoreConfigs: []talismanrc.IgnoreConfig{
				&talismanrc.FileIgnoreConfig{
					FileName: "some.txt",
					Checksum: "sha1",
				},
			},
		}
		ie := IgnoreEvaluator{calculator: checksumCalculator, talismanRC: &ignoreConfig}
		addition := gitrepo.Addition{Name: "some.txt", Path: "some.txt"}
		checksumCalculator.EXPECT().CalculateCollectiveChecksumForPattern("some.txt").Return("sha1")

		required := ie.isScanNotRequired(addition)

		assert.True(t, required)
	})

}

type sillyChecksumCalculator struct{}

func (scc *sillyChecksumCalculator) CalculateCollectiveChecksumForPattern(fileNamePattern string) string {
	return "silly"
}
func (scc *sillyChecksumCalculator) SuggestTalismanRC(fileNamePatterns []string) string {
	return ""
}

func TestDeterminingFilesToIgnore(t *testing.T) {
	tRC := talismanrc.TalismanRC{
		IgnoreConfigs: []talismanrc.IgnoreConfig{
			&talismanrc.FileIgnoreConfig{
				FileName: "some.txt",
				Checksum: "silly",
			},
			&talismanrc.FileIgnoreConfig{
				FileName: "other.txt",
				Checksum: "serious",
			},
			&talismanrc.FileIgnoreConfig{
				FileName:        "ignore-contents",
				IgnoreDetectors: []string{"filecontent"},
			},
		},
	}
	ie := IgnoreEvaluator{&sillyChecksumCalculator{}, &tRC}

	t.Run("Should ignore file based on checksum", func(t *testing.T) {
		assert.True(t, ie.ShouldIgnore(gitrepo.Addition{Path: "some.txt"}, ""))
	})

	t.Run("Should not ignore file if checksum doesn't match", func(t *testing.T) {
		assert.False(t, ie.ShouldIgnore(gitrepo.Addition{Path: "other.txt"}, ""))
	})

	t.Run("Should ignore if detector is disabled for file", func(t *testing.T) {
		assert.True(t, ie.ShouldIgnore(gitrepo.Addition{Path: "ignore-contents"}, "filecontent"))
	})

	t.Run("Should not ignore if a different detector is disabled for file", func(t *testing.T) {
		assert.False(t, ie.ShouldIgnore(gitrepo.Addition{Path: "ignore-contents"}, "filename"))
	})
}
