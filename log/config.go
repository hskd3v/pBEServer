package log

import "os"

const (
	productionMode  = "production"
	developmentMode = "development"
)

const (
	envLogLevel = "BE_LOG_LEVEL"
	envLogMode  = "BE_LOG_MODE"

	defLogLevel = "debug"
	defLogMode  = developmentMode
)

var (
	cfgLogLevel = EnvStr(envLogLevel, defLogLevel)
	cfgLogMode  = EnvStr(defLogMode, defLogMode)
)

// IsProduction return true if it is on production environment
func IsProduction() bool {
	return cfgLogMode == productionMode
}

// EnvStr func
func EnvStr(aKey string, aDef string) string {
	_val, _ok := os.LookupEnv(aKey)
	if _ok {
		return _val
	}
	return aDef
}
