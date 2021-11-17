package util

import "github.com/rs/zerolog"

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
