package logger

import (
	"bytes"
	"encoding/json"
	"sync"
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

type logHole struct {
	sync.Mutex
	hole map[string][]map[string]string
}

func newHole() *logHole {
	return &logHole{
		Mutex: sync.Mutex{},
		hole:  map[string][]map[string]string{},
	}
}

func (l *logHole) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()
	var objmap map[string]string
	json.Unmarshal(p, &objmap)

	l.hole[objmap["level"]] = append(l.hole[objmap["level"]], objmap)

	return len(p), nil
}

func (l *logHole) getLastLog(level string) string {
	l.Lock()
	defer l.Unlock()
	count := len(l.hole[level])
	if count == 0 {
		return ""
	}
	return l.hole[level][count-1]["message"]
}
func (l *logHole) getLogs(level string) []string {
	var messages []string

	for i := range l.hole[level] {
		messages = append(messages, l.hole[level][i]["message"])
	}

	return messages
}

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
	hole := newHole()
	params := Params{
		LogLevel: "aasdasd",
		output:   hole,
	}
	_ = New(params)
	require.Equal(t, []string{"failed to parse log level", "this vars didn't set: [Version Component]"}, hole.getLogs("error"))
}

func TestNewLogger_VersionNotSet_ReturnError(t *testing.T) {
	hole := newHole()
	params := Params{
		LogLevel:  "error",
		Component: "test",
		output:    hole,
	}

	_ = New(params)

	require.Equal(t, "this vars didn't set: [Version]", hole.getLastLog("error"))
}

func TestNewLogger_ComponentNotSet_ReturnError(t *testing.T) {
	hole := newHole()

	params := Params{
		LogLevel:   "error",
		AppVersion: "1.0",
		output:     hole,
	}

	_ = New(params)
	require.Equal(t, hole.getLastLog("error"), "this vars didn't set: [Component]")
}

func TestNewLogger_CheckLogLevel_IfSetErrorLevel_DontPrintInfoLogs(t *testing.T) {
	hole := newHole()

	params := Params{
		LogLevel:   "error",
		Component:  "test",
		AppVersion: "1.0",
		output:     hole,
	}

	logger := New(params)
	logger.Error().Msg("test")

	require.Equal(t, "", hole.getLastLog("info"))
}
