package config

import (
	"log"
	"os"
)

func GetToken() string {
	return getEnvVariable("TOKEN")
}

func GetNatsURI() string {
	return getEnvVariable("NATSURI")
}

func GetCeduladosServiceURL() string {
	return getEnvVariable("CEDULADOS_SERVICE_URL")
}

func GetUserServiceURL() string {
	return getEnvVariable("TGBUSER_SERVICE_URL")
}

func GetEnvironment() string {
	return getEnvVariable("ENV")
}

func getEnvVariable(key string) string {
	if os.Getenv(key) == "" {
		log.Fatal("error getting key", key)
	}
	return os.Getenv(key)
}
