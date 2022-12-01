package utils_env

import (
	"log"
	"os"
	"strings"
	"teniditter-server/cmd/global/console"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("APP_MODE") == "prod" {
		return
	}

	// Default exec scope
	envPath := "../../local.env"

	// Find root of project
	execPath, err := os.Getwd()
	if err == nil {
		var n int

		paths := strings.Split(execPath, "/")
		for i, path := range paths {
			if path == "server" {
				n = len(paths[i+1:])
			}
		}

		envPath = strings.Repeat("../", n) + "local.env"
	}

	log.Println(envPath)
	if err = godotenv.Load(envPath); err != nil {
		log.Fatalf("ABORTED: Cannot Load env file. Err: %s", err)
	}
	console.Log("[dev]: .env file loaded")
}
