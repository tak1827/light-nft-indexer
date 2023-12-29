package util

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func LoadDotenv(filepaths ...string) error {
	if err := godotenv.Load(filepaths...); err != nil {
		return errors.Wrap(err, "failed loading env file")
	}
	return nil
}

func IsNoEnvErr(err error) bool {
	return strings.Contains(err.Error(), "no such file or directory")
}
