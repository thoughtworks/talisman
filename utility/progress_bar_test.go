package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProgressBar(t *testing.T) {
	t.Run("Start should start inner progress bar", func(t *testing.T) {
		progressBar := &defaultProgressBar{}

		progressBar.Start(10)

		assert.True(t, progressBar.bar.IsStarted())
	})

	t.Run("Increment should update inner progress bar", func(t *testing.T) {
		progressBar := &defaultProgressBar{}

		progressBar.Start(10)
		progressBar.Increment()

		assert.Equal(t, int64(1), progressBar.bar.Current())
	})

	t.Run("Finish should finish progress bar", func(t *testing.T) {
		progressBar := &defaultProgressBar{}

		progressBar.Start(2)
		progressBar.Increment()
		progressBar.Increment()
		progressBar.Finish()

		assert.False(t, progressBar.bar.IsStarted())
	})

}
