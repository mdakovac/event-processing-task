package env_vars

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariablesType struct {
	PUBSUB_EMULATOR_HOST string
	PUBSUB_PROJECT_ID    string
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
}
