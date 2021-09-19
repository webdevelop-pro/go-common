package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type LogRaw struct {
	Level    string `json:"level"`
	Severity string `json:"severity"`
	Time     int    `json:"time"`
	Message  string `json:"message"`
}

var testMsg = "test message"
var defaultLogger = GetDefaultLogger(nil)

func TestGetDefaultLogger_checkLevelIsSet(t *testing.T) {
	var buf bytes.Buffer

	logger := GetDefaultLogger(&buf)

	logger.Error().Msg(testMsg)

	var resultLogStr LogRaw

	json.Unmarshal([]byte(buf.String()), &resultLogStr)

	require.Equal(t, "error", resultLogStr.Level)
}

func TestGetDefaultLogger_checkSeverityIsSet(t *testing.T) {
	var buf bytes.Buffer

	logger := GetDefaultLogger(&buf)

	logger.Error().Msg(testMsg)

	var resultLogStr LogRaw

	json.Unmarshal([]byte(buf.String()), &resultLogStr)

	require.Equal(t, "ERROR", resultLogStr.Severity)
}

func TestGetDefaultLogger_checkMessageIsSet(t *testing.T) {
	var buf bytes.Buffer

	logger := GetDefaultLogger(&buf)

	logger.Error().Msg(testMsg)

	var resultLogStr LogRaw

	json.Unmarshal([]byte(buf.String()), &resultLogStr)

	require.Equal(t, testMsg, resultLogStr.Message)
}

func TestGetDefaultLogger_checkTimeIsSet(t *testing.T) {
	var buf bytes.Buffer

	logger := GetDefaultLogger(&buf)

	logger.Error().Msg(testMsg)

	var resultLogStr LogRaw

	json.Unmarshal([]byte(buf.String()), &resultLogStr)

	require.True(t, time.Unix(int64(resultLogStr.Time), 0).Before(time.Now()))
	require.True(t, time.Unix(int64(resultLogStr.Time), 0).After(time.Now().Add(time.Minute*-1)))
}

func TestGetDefaultLogger_outputIsNil(t *testing.T) {
	require.NotPanics(t, func() {
		logger := GetDefaultLogger(nil)
		logger.Error().Msg(testMsg)
	})
}

func TestNewLogger_ParseLevelError_ReturnError(t *testing.T) {
	params := Params{
		LogLevel: "aasdasd",
	}

	logger, err := NewLogger(params)

	require.NotNil(t, err)
	require.Equal(t, logger, defaultLogger)
}

func TestNewLogger_VersionNotSet_ReturnError(t *testing.T) {
	params := Params{
		LogLevel:  "error",
		Component: "test",
	}

	logger, err := NewLogger(params)

	require.NotNil(t, err)
	require.Equal(t, logger, defaultLogger)
	require.Equal(t, err.Error(), "this vars didn't set: [Version]")
}

func TestNewLogger_ComponentNotSet_ReturnError(t *testing.T) {
	params := Params{
		LogLevel:   "error",
		AppVersion: "1.0",
	}

	logger, err := NewLogger(params)

	require.NotNil(t, err)
	require.Equal(t, logger, defaultLogger)
	require.Equal(t, err.Error(), "this vars didn't set: [Component]")
}

func TestNewLogger_CheckLogLevel_IfSetErrorLevel_DontPrintInfoLogs(t *testing.T) {
	var buf bytes.Buffer
	params := Params{
		LogLevel:   "error",
		Component:  "test",
		AppVersion: "1.0",
		output:     &buf,
	}

	logger, err := NewLogger(params)

	logger.Info().Msg("test")

	require.Nil(t, err)
	require.Equal(t, "", buf.String())
}
