package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't load .env file")
	}
}
