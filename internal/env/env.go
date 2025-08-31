package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	return val
}

func GetInt(key string, defaultVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}

func GetBool(key string, defaultVal bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}

	return boolVal
}

func GetDuration(key string, defaultVal time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	durationVal, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}

	return durationVal
}
