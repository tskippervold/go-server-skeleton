package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var connection *sqlx.DB

func GetConnection() *sqlx.DB {
	return connection
}

/*
	Consider using pgx and sqlx instead of gorm ORM.
	https://github.com/jackc/pgx
	https://github.com/jmoiron/sqlx
	https://github.com/jackc/pgx/issues/81
*/

func Connect(host string, port uint16, name string, user string, pass string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, name)
	db, err := sqlx.Connect("pgx", connStr)

	if err != nil {
		return err
	}

	connection = db
	return nil
}

func ConnectToTest() error {
	db, err := sqlx.Connect("sqlite3", "sample.db")

	if err != nil {
		return err
	}

	connection = db
	return nil
}

func Migrate() error {
	driver, err := postgres.WithInstance(connection.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://../internal/db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
