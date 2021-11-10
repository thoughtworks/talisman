package main

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_setLogLevel(t *testing.T) {
	levels := []string{"error", "warn", "info", "debug", "unknown"}
	expectedLogrusLevels := []logrus.Level{
		logrus.ErrorLevel, logrus.WarnLevel,
		logrus.InfoLevel, logrus.DebugLevel, logrus.ErrorLevel}

	for idx, level := range levels {
		options.LogLevel = level
		setLogLevel()
		assert.True(
			t,
			logrus.IsLevelEnabled(expectedLogrusLevels[idx]),
			fmt.Sprintf("expected level to be %v for options.LogLevel = %s", expectedLogrusLevels[idx], level))

		options.Debug = true
		setLogLevel()
		assert.True(
			t,
			logrus.IsLevelEnabled(logrus.DebugLevel),
			"expected level to be debug when options.Debug is set")
	}
}
