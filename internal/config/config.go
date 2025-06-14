package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/prajwalbharadwajbm/broker/internal/utils"
)

type GeneralConfig struct {
	Env      string
	LogLevel string
	Port     int
}

type appConfig struct {
	GeneralConfig GeneralConfig
	DB            DB
	JWTSecret     string
}

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

func LoadConfigs() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env files: ", err)
	}

	loadGeneralCongigs()
	loadDatabaseConfigs()
	loadJWTConfigs()
}

var AppConfigInstance appConfig

func loadGeneralCongigs() {
	AppConfigInstance.GeneralConfig.Env = utils.GetEnv("APP_DEV", "dev")
	AppConfigInstance.GeneralConfig.LogLevel = utils.GetEnv("LOG_LEVEL", "info")
	AppConfigInstance.GeneralConfig.Port = utils.GetEnv("PORT", 8080)
}

func loadDatabaseConfigs() {
	AppConfigInstance.DB.Host = utils.GetEnv("DB_HOST", "localhost")
	AppConfigInstance.DB.Port = utils.GetEnv("DB_PORT", 5432)
	AppConfigInstance.DB.User = utils.GetEnv("POSTGRES_USER", "postgres")
	AppConfigInstance.DB.Password = utils.GetEnv("POSTGRES_PASSWORD", "")
	AppConfigInstance.DB.DBname = utils.GetEnv("POSTGRES_DB", "broker-platform")
}

func loadJWTConfigs() {
	AppConfigInstance.JWTSecret = utils.GetEnv("JWT_SECRET", "")
}
