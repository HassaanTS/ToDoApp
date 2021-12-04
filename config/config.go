package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	MongoHost         string
	MongoPort         string
	MongoDBName       string
	MongoDBCollection string
}

var GlobalConfig Config

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("couldn't load environment from .env file...%w", err)
	}

	GlobalConfig.AppPort = os.Getenv("APP_PORT")
	GlobalConfig.MongoHost = os.Getenv("MONGODB_HOSTNAME")
	GlobalConfig.MongoPort = os.Getenv("MONGODB_PORT")
	GlobalConfig.MongoDBName = os.Getenv("MONGODB_NAME")
	GlobalConfig.MongoDBCollection = os.Getenv("MONGODB_COLLECTION")

	return nil
}
