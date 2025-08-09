package env

import (
	"os"
	"strconv"
)

func ParseUIntEnv(name string, defValue *uint) {
	if envValue := os.Getenv(name); envValue != "" {
		if value, err := strconv.ParseUint(envValue, 10, 32); err == nil {
			*defValue = uint(value)
		}
	}
}
