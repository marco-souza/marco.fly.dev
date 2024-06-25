package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Ctx     context.Context = context.Background()
	client  *sql.DB
	Queries *sqlc.Queries
)

//go:embed schema.sql
var ddl string

func Init(file string) error {
	log.Println("init db")

	// TODO: use a file
	if file == "" {
		file = ":memory:"
	}

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	// setup client and context
	client = db

	// setup the database schema
	if _, err := db.ExecContext(Ctx, ddl); err != nil {
		log.Println("error configuring db tables: ", err)
	}

	Queries = sqlc.New(db)
	return nil
}

func Close() error {
	log.Println("closing db")

	Ctx = nil
	Queries = nil
	return client.Close()
}
