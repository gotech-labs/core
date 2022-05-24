package env

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gotech-labs/core/log"
)

// SetenvIfNotExists is ...
func SetenvIfNotExists(key, value string) {
	if os.Getenv(key) == "" {
		if err := os.Setenv(key, value); err != nil {
			log.Panic().Err(err).Msgf("failed to set env [key=%s, value=%s]", key, value)
		}
	}
}

// GetString is ...
func GetString(key, fallback string) string {
	v := os.Getenv(key)
	if len(v) > 0 {
		return v
	}
	return fallback
}

// GetInt is ...
func GetInt(key string, fallback int) int {
	value := GetString(key, "")
	if len(value) > 0 {
		if num, err := strconv.Atoi(value); err == nil {
			return num
		}
	}
	return fallback
}

// GetBool is ...
func GetBool(key string, fallback bool) bool {
	value := GetString(key, "")
	switch strings.ToLower(value) {
	case "true":
		return true
	default:
		return fallback
	}
}

// GetDuration is ...
func GetDuration(key string, fallback time.Duration) (time.Duration, error) {
	v := GetString(key, "")
	if len(v) > 0 {
		return time.ParseDuration(v)
	}
	return fallback, nil
}

// MustGetDuration is ...
func MustGetDuration(key string, fallback time.Duration) time.Duration {
	v, err := GetDuration(key, fallback)
	if err != nil {
		log.Panic().Err(err).Msgf("failed to parse duration [key=%s, value=%s]", key, os.Getenv(key))
	}
	return v
}
