package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import the file source driver
	"log"
	"trading-ace/src/config"
)

func MigrateDB(migrationSource string, dbConfig *config.DatabaseConfig) {
	migrations, err := migrate.New(migrationSource, dbConfig.GetUrl())

	if err != nil {
		log.Fatal(err)
	}

	err = migrations.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	if err == nil {
		log.Println("Database migration successful")
	}
}
