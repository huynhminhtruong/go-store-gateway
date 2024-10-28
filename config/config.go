package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV") // Possible values for development/production
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL") // Database connection URL
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT") // Book Client service port
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetBookServiceUrl() string {
	return getEnvironmentValue("BOOK_SERVICE_URL") // This will be in the Book service env params
}

func getEnvironmentValue(key string) string { // Validates env param exists and gets it
	if os.Getenv(key) == "" { // GetEnv returns the string
		log.Fatalf("%s environment variable is missing", key)
	}

	return os.Getenv(key)
}
