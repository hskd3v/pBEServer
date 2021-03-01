package utils

import (
	"os"
	"strconv"
)

// EnvStr func
func EnvStr(aKey string, aDef string) string {
	_val, _ok := os.LookupEnv(aKey)
	if _ok {
		return _val
	}
	return aDef
}

// EnvInt func
func EnvInt(aKey string, aDef int) int {
	_val, _ok := os.LookupEnv(aKey)
	if _ok {
		_i, err := strconv.Atoi(_val)
		if err == nil {
			return _i
		}
	}
	return aDef
}

// EnvInt64 func
func EnvInt64(aKey string, aDef int64) int64 {
	_val, _ok := os.LookupEnv(aKey)
	if _ok {
		_i, err := strconv.ParseInt(_val, 10, 64)
		if err == nil {
			return _i
		}
	}
	return aDef
}
