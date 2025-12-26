package config

import (
	"fmt"
	"os"
	"strings"
)

func RequireEnv(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", fmt.Errorf("missing required env var %s", key)
	}
	return v, nil
}

func PreviewBaseDomain() (string, error) {
	return RequireEnv("PREVIEW_BASE_DOMAIN")
}
