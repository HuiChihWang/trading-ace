package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	return appConfig
}

func LoadConfig() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return err
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Error unmarshalling config data: %v", err)
		return err
	}

	appConfig.AppEnv = env

	return nil
}
