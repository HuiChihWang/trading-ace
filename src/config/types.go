package config

import (
	"fmt"
	"time"
)

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

type RedisConnectionConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database int    `mapstructure:"db"`
}
type RedisConfig struct {
	Job *RedisConnectionConfig `mapstructure:"job"`
}

type EthereumNodeConfig struct {
	SocketUrl string `mapstructure:"socket"`
}

type CampaignConfig struct {
	CampaignStartTime string `mapstructure:"start_time"`
	Weeks             int    `mapstructure:"weeks"`
}

func (c *CampaignConfig) GetCampaignStartTime() time.Time {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, c.CampaignStartTime)
	return t
}

type AppConfig struct {
	AppEnv       string
	Database     *DatabaseConfig     `mapstructure:"database"`
	EthereumNode *EthereumNodeConfig `mapstructure:"ethereum_node"`
	Campaign     *CampaignConfig     `mapstructure:"campaign"`
	Redis        *RedisConfig        `mapstructure:"redis"`
}
