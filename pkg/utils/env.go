package utils

import (
	"log"
	"os"
)

const (
	envPrefix = "ZNEIXBOT_"
)

func GetEnv(envName string, isRequired bool) (value string, exists bool) {
	value, envExists := os.LookupEnv(envPrefix + envName)

	if !envExists && isRequired {
		log.Fatalf("Missing required %s environment variable", envPrefix+envName)
	}

	return value, envExists
}
