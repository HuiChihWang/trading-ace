package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"trading-ace/src/config"
)

var db *sql.DB

func GetDBInstance() *sql.DB {
	return db
}

func CreateDBInstance(databaseConfig *config.DatabaseConfig) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.DBName,
		databaseConfig.Username,
		databaseConfig.Password,
	)

	log.Println("Connecting to database... " + connStr)

	var err error
	newDb, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = newDb.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return newDb
}

func InitDatabase() error {
	dbConfig := config.GetAppConfig().Database
	db = CreateDBInstance(&dbConfig)
	return nil
}
