package config

import "fmt"

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func (d *DatabaseConfig) GetUrl() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", d.Driver, d.Username, d.Password, d.Host, d.Port, d.DBName)
}

func (d *DatabaseConfig) GetConnectionStr() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", d.Host, d.Port, d.DBName, d.Username, d.Password)
}

type EthereumNodeConfig struct {
	SocketUrl string `mapstructure:"socket"`
}

type AppConfig struct {
	AppEnv       string
	Database     *DatabaseConfig     `mapstructure:"database"`
	EthereumNode *EthereumNodeConfig `mapstructure:"ethereum_node"`
}
