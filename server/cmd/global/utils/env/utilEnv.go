package utils_env

import (
	"log"
	"os"
	"teniditter-server/cmd/global/console"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("APP_MODE") == "prod" {
		return
	}

	err := godotenv.Load("../../local.env")
	if err != nil {
		log.Fatalf("ABORTED: Cannot Load env file. Err: %s", err)
	}
	console.Log("[dev]: .env file loaded", console.Info)
}
