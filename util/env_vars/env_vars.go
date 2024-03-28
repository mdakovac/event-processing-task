package env_vars

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariablesType struct {
	PUBSUB_EMULATOR_HOST string
}

var EnvVariables EnvVariablesType

func SetEnvVars() {
	log.Println("SETTING ENV VARS")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	EnvVariables = EnvVariablesType{}

	EnvVariables.PUBSUB_EMULATOR_HOST = os.Getenv("PUBSUB_EMULATOR_HOST")
}
