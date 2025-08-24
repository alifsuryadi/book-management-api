package database

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(databaseURL string) error {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFiles,
		Root:       "migrations",
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}