package env_vars

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariablesType struct {
	PUBSUB_EMULATOR_HOST   string
	PUBSUB_PROJECT_ID      string
	EXCHANGE_RATES_API_URL string
	DB_CONNECTION_URL      string
}

var EnvVariables EnvVariablesType

func SetEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	EnvVariables = EnvVariablesType{}

	EnvVariables.PUBSUB_EMULATOR_HOST = os.Getenv("PUBSUB_EMULATOR_HOST")
	EnvVariables.PUBSUB_PROJECT_ID = os.Getenv("PUBSUB_PROJECT_ID")
	EnvVariables.EXCHANGE_RATES_API_URL = os.Getenv("EXCHANGE_RATES_API_URL")
	EnvVariables.DB_CONNECTION_URL = os.Getenv("DB_CONNECTION_URL")
}
