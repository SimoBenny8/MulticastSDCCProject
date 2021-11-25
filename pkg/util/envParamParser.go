package util

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*This functions parse the environmental variables of Docker*/

//GetEnvStr return value of env variable
func GetEnvString(key string) (string, error) {
	isValid := ValidateKeyName(key)
	if isValid == false {
		return "", errors.New("invalid variable name")
	}

	v := strings.TrimSpace(os.Getenv(key))
	if len(v) == 0 {
		return v, errors.New("environment variable empty")
	}

	return v, nil
}

// GetEnvInt return value of env variable.
func GetEnvInt(key string) (int, error) {
	s, err := GetEnvString(key)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return v, err
}

// GetEnvBool return value of env variable.
func GetEnvBool(key string) (bool, error) {
	s, err := GetEnvString(key)
	if err != nil {
		return false, err
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}

	return v, err
}

// GetEnvStringWithDefault return value of env variable if any or return default value if error
func GetEnvStringWithDefault(key string, defaultValue string) string {
	v, err := GetEnvString(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetEnvIntWithDefault return value of env variable if any or return default value if error
func GetEnvIntWithDefault(key string, defaultValue int) int {
	v, err := GetEnvInt(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetEnvBoolWithDefault return value of env variable if any or return default value if error
func GetEnvBoolWithDefault(key string, defaultValue bool) bool {
	v, err := GetEnvBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

var pattern = "^[a-zA-Z_]+[a-zA-Z0-9_]*$"
var IsLetter = regexp.MustCompile(pattern).MatchString

// ValidateKeyName must ensure env variable name only contain alphanumeric + underscore
func ValidateKeyName(key string) bool {
	result := IsLetter(key)
	return result
}
