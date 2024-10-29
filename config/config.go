package config

import (
	"fmt"
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
	portStr := getEnvironmentValue("BOOK_PORT") // Book Client service port
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetServiceURL(ipSrv, portSrv string) string {
	ip := getEnvironmentValue(ipSrv)
	port := getEnvironmentValue(portSrv)
	return fmt.Sprintf("%s:%s", ip, port)
}

func getEnvironmentValue(key string) string { // Validates env param exists and gets it
	if os.Getenv(key) == "" { // GetEnv returns the string
		log.Fatalf("%s environment variable is missing", key)
	}

	return os.Getenv(key)
}
