package util

import (
	"errors"
	"github.com/rs/zerolog"
	"os"
)

func CheckError(err error, log *zerolog.Logger) {
	log.Fatal().Err(err)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetAWSRegionEnvVar() (string, error) {
	var err error
	val, present := os.LookupEnv("AWS_DEFAULT_REGION")
	if !present {
		return "", errors.New("AWS_DEFAULT_REGION is not set")
	}
	if val == "" {
		return "", errors.New("AWS_DEFAULT_REGION is set to an empty string")
	}
	return val, err
}
