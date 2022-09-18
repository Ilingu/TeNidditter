package main

import (
	"log"
	"os"

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
}
