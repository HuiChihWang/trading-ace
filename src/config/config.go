package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

var appConfig *AppConfig
var once sync.Once

func GetAppConfig() *AppConfig {
	once.Do(func() {
		loadedConfig, err := LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		appConfig = loadedConfig
	})
	return appConfig
}

func LoadConfig() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("json")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var loadedConfig AppConfig
	if err := viper.Unmarshal(&loadedConfig); err != nil {
		log.Fatalf("Error unmarshalling config data: %v", err)
		return nil, err
	}
	loadedConfig.AppEnv = env

	return &loadedConfig, nil
}
