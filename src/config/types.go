package config

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type EthereumNodeConfig struct {
	SocketUrl string `mapstructure:"socket"`
}

type AppConfig struct {
	AppEnv       string
	Database     DatabaseConfig     `mapstructure:"database"`
	EthereumNode EthereumNodeConfig `mapstructure:"ethereum_node"`
}
