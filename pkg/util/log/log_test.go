package log

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestCustomFormatter_Format(t *testing.T) {

	logger := logrus.New()
	f := &CustomFormatter{}
	data := map[string]interface{}{
		"a":            "1",
		"b":            2,
		"time":         time.Now().String(),
		"msg":          "ggg",
		"func":         "111",
		"file":         "1",
		"logrus_error": "1",
		"level":        DebugLevel,
	}
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Data:    data,
		Time:    time.Time{},
		Level:   DebugLevel,
		Caller:  &runtime.Frame{},
		Message: "ggg",
		Buffer:  nil,
		Context: nil,
	}
	entry.Caller = &runtime.Frame{}

	f.ForceColors = true
	f.EnvironmentOverrideColors = true
	os.Setenv("CLICOLOR_FORCE", "0")
	f.Format(entry)
	logger.SetFormatter(f)

	os.Unsetenv("CLICOLOR_FORCE")
	os.Setenv("CLICOLOR", "0")
	f.Format(entry)
	logger.SetFormatter(f)

	f.ForceColors = true
	f.EnvironmentOverrideColors = true
	os.Setenv("CLICOLOR_FORCE", "true")
	logger.SetFormatter(f)

	// test level
	//f.PadLevelText=true
	entry.Level = TraceLevel
	f.DisableLevelTruncation = false
	f.PadLevelText = true
	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	entry.Level = InfoLevel
	f.DisableLevelTruncation = false
	f.PadLevelText = false
	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	entry.Level = WarnLevel
	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	entry.Level = ErrorLevel
	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	// test DisableTimestamp
	f.DisableTimestamp = true
	logger.SetFormatter(f)

	f.DisableTimestamp = false
	f.FullTimestamp = false
	logger.SetFormatter(f)
	//
	f.DisableTimestamp = true
	f.FullTimestamp = true
	logger.SetFormatter(f)

	f.DisableTimestamp = true
	f.FullTimestamp = false
	logger.SetFormatter(f)

	// test entry HasCaller
	entry.Logger = logrus.New()
	entry.Logger.ReportCaller = true
	entry.Caller = &runtime.Frame{}
	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	// test callerpretty
	entry.Logger = logrus.New()
	entry.Logger.ReportCaller = true
	fnc := func(*runtime.Frame) (function string, file string) {
		return "", ""
	}
	f.CallerPrettyfier = fnc

	f.Format(entry)
	logger.SetFormatter(f)
	logger.Info("dddd")

	// test f.SortingFunc
	f.SortingFunc = func(strings []string) {

	}
	logger.SetFormatter(f)
	logger.Info("dddd")

	f.ForceColors = false
	logger.SetFormatter(f)

	// ------------------------ //
	// test f.DisableSorting
	logger2 := logrus.New()
	f2 := &CustomFormatter{}
	data2 := map[string]interface{}{
		"a":            "1",
		"b":            2,
		"time":         time.Now().String(),
		"msg":          "ggg",
		"func":         "111",
		"file":         "1",
		"logrus_error": "1",
		"level":        DebugLevel,
	}
	entry2 := &logrus.Entry{
		Logger:  logrus.New(),
		Data:    data2,
		Time:    time.Time{},
		Level:   DebugLevel,
		Caller:  &runtime.Frame{},
		Message: "ggg",
		Buffer:  nil,
		Context: nil,
	}
	entry2.Caller = &runtime.Frame{}

	f2.DisableSorting = false
	f2.CallerPrettyfier = func(*runtime.Frame) (function string, file string) {
		return "ddd", "ddd"
	}

	f2.Format(entry2)
	logger2.SetFormatter(f2)

}

func TestNeedsQuoting(t *testing.T) {
	f := &CustomFormatter{}
	f.QuoteEmptyFields = true
	// test return true
	assert.Equal(t, true, f.needsQuoting(""))

	// test return true
	assert.Equal(t, true, f.needsQuoting(","))

	// test return false
	assert.Equal(t, false, f.needsQuoting("-"))
}

//func (f *CustomFormatter) needsQuoting(text string) bool {
//	if f.QuoteEmptyFields && len(text) == 0 {
//		return true
//	}
//	for _, ch := range text {
//		if !((ch >= 'a' && ch <= 'z') ||
//			(ch >= 'A' && ch <= 'Z') ||
//			(ch >= '0' && ch <= '9') ||
//			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
//			return true
//		}
//	}
//	return false
//}
