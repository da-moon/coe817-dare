package log

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var stdLog *defaultLogger

// SetTestLogger - initializes logger for tests
func SetTestLogger(t *testing.T) {
	stdLog = newTestLogger(t)
}

// SetDefaultLogger ...
func SetDefaultLogger() {
	stdLog = newDefaultLogger()
}

// Debug log with level debug
func Debug(msg string) {
	if stdLog == nil {
		return
	}
	stdLog.Debug(msg, nil)
}

// Info log with level Info
func Info(msg string) {
	if stdLog == nil {
		return
	}
	stdLog.Info(msg, nil)
}

// Warn log with level Warn
func Warn(msg string) {
	if stdLog == nil {
		return
	}
	stdLog.Warn(msg, nil)
}

// Error log with level Error
func Error(msg string) {
	if stdLog == nil {
		return
	}
	stdLog.Error(msg, nil)
}

// defaultLogger - Default implementation of Logger
type defaultLogger struct {
	Log *log.Logger
}

// newDefaultLogger - Create new default logger
func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		Log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// newDefaultLogger - Create new default logger
func newTestLogger(t *testing.T) *defaultLogger {
	return &defaultLogger{
		Log: log.New(testWriter{t}, "test", log.LstdFlags),
	}
}
func (d *defaultLogger) addLogLevel(level, msg string) string {
	return fmt.Sprintf("[%s] %s", level, msg)
}

func (d *defaultLogger) addExtraFields(extraFields map[string]interface{}, msg string) string {
	extraString := ""
	for k, v := range extraFields {
		extraString = fmt.Sprintf("%s%s=%v ", extraString, k, v)
	}
	if extraString != "" {
		msg = extraString + msg
	}
	return msg
}

// Debug -
func (d *defaultLogger) Debug(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("DEBUG", msg)
	d.Log.Print(msg)
}

// Info -
func (d *defaultLogger) Info(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("INFO", msg)
	d.Log.Print(msg)
}

// Warn -
func (d *defaultLogger) Warn(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("WARN", msg)
	d.Log.Print(msg)
}

// Error -
func (d *defaultLogger) Error(msg string, extraFields map[string]interface{}) {
	msg = d.addExtraFields(extraFields, msg)
	msg = d.addLogLevel("ERROR", msg)
	d.Log.Print(msg)
}

type testWriter struct {
	t *testing.T
}

// Write ...
func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}
