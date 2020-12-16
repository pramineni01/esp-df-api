package config

import (
	"os"
)

// Settings
// ----------------------------------------------

func ENDPOINT_DF_API_INTERNAL() string {
	if val := os.Getenv("ENDPOINT_DF_API_INTERNAL"); val != "" {
		return val
	}
	return DEFAULT_ENDPOINT_DF_API_INTERNAL
}
