package utils

import (
	"os"
	"strconv"
)

// GetEnv returns the environment variable value if it exists, otherwise returns the fallback value.
// The return type matches the type of the fallback parameter.
func GetEnv[T any](key string, fallback T) T {
	if value, exists := os.LookupEnv(key); exists {
		switch any(fallback).(type) {
		case string:
			return any(value).(T)
		case int:
			if intVal, err := strconv.Atoi(value); err == nil {
				return any(intVal).(T)
			}
		case bool:
			if boolVal, err := strconv.ParseBool(value); err == nil {
				return any(boolVal).(T)
			}
		case float64:
			if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				return any(floatVal).(T)
			}
		}
	}
	return fallback
}
