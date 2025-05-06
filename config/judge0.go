package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Judge0Config struct {
	Key  string
	Host string
}

func LoadJudge0Config() Judge0Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: failed to load .env file")
	}

	return Judge0Config{
		Key:  os.Getenv("RAPIDAPI_KEY"),
		Host: os.Getenv("RAPIDAPI_HOST"),
	}
}
