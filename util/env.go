package util

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func LoadDotenv(filepaths ...string) error {
	if err := godotenv.Load(filepaths...); err != nil {
		return fmt.Errorf("failed loading env file: %w", err)
	}
	return nil
}

func IsNoEnvErr(err error) bool {
	return strings.Contains(err.Error(), "no such file or directory")
}
