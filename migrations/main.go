package main

import (
	"trading-ace/src/config"
	"trading-ace/src/database"
)

func main() {
	dbConfig := config.GetAppConfig().Database
	database.MigrateDB("file://migrations", dbConfig)
}
