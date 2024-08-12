package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"trading-ace/src/config"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDBInstance() *sql.DB {
	once.Do(func() {
		dbConfig := config.GetAppConfig().Database
		db = CreateDBInstance(dbConfig)
		log.Println("Database connection established: ", dbConfig.GetUrl())
	})
	return db
}

func CreateDBInstance(databaseConfig *config.DatabaseConfig) *sql.DB {
	newDb, err := sql.Open(databaseConfig.Driver, databaseConfig.GetConnectionStr())

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
