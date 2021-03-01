package utils

import (
	"bufio"
	"io"

	"github.com/go-openapi/strfmt"
)

// IsArrayFromReader returns if the string is an json array
func IsArrayFromReader(aReader io.Reader) bool {
	_isArray := false
	_reader := bufio.NewReader(aReader)
	for _c, _ := _reader.ReadByte(); ; _c, _ = _reader.ReadByte() {
		if _c == ' ' || _c == '\t' || _c == '\r' || _c == '\n' {
			continue
		}
		_isArray = _c == '['
		break
	}
	return _isArray
}

// IsArray returns if the string is an json array
func IsArray(aReader *[]byte) bool {
	_isArray := false
	for _, c := range *aReader {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		_isArray = c == '['
		break
	}
	return _isArray
}

// IsValidUUID validates a UUID
func IsValidUUID(aUUID string) bool {
	return strfmt.IsUUID(aUUID)
	//r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	//return r.MatchString(uuid)
}

// UUIDHasCurlyBraces return if a string has brackets
func UUIDHasCurlyBraces(aString string) bool {
	_hasBarckets := false

	if len(aString) > 0 {
		if aString[0] == '{' {
			_hasBarckets = true
		}
	}
	return _hasBarckets
}

// UUIDAddBrackets adds brackets
func UUIDAddCurlyBraces(aString string) string {
	if UUIDHasCurlyBraces(aString) {
		return aString
	}
	return "{" + aString + "}"
}

// UUIDRemoveCurlyBraces removes brackets
func UUIDRemoveCurlyBraces(aString string) string {
	if !UUIDHasCurlyBraces(aString) {
		return aString
	}
	_len := len(aString)
	if aString[_len-1] == '}' {
		_len--
	}

	return aString[1:_len]
}
