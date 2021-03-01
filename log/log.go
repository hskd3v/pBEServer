package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Log is the global variable for the customized log
var Log *logrus.Logger

// Init initializes log
func Init() {
	_level, _error := logrus.ParseLevel(cfgLogLevel)
	if _error != nil {
		_level = logrus.DebugLevel
	}

	Log = &logrus.Logger{
		Level: _level,
		Out:   os.Stdout,
	}

	if IsProduction() {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}
}

// ParseFields attempts to convert a string to a field. eg.: "f:v" to {"f","v"}
func ParseFields(aTags ...string) logrus.Fields {
	_result := make(logrus.Fields, len(aTags))
	for _, _tag := range aTags {
		_els := strings.Split(_tag, ":")
		_result[strings.TrimSpace(_els[0])] = strings.TrimSpace(_els[1])
	}
	return _result
}
