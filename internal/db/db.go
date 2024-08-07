package db

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"
	"github.com/marco-souza/marco.fly.dev/internal/di"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var logger = slog.With("service", "db")

//go:embed schema.sql
var dbSchema string

type DatabaseService struct {
	File    string
	Ctx     context.Context
	client  *sql.DB
	Queries *sqlc.Queries
}

func New() *DatabaseService {
	cfg, err := di.Inject(config.Config{})
	if err != nil {
		panic(err)
	}

	file := cfg.SqliteUrl
	if file == "" {
		logger.Info("use in-memory db")
		file = ":memory:"
	}

	ctx := context.Background()
	client, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	queries := sqlc.New(client)
	return &DatabaseService{file, ctx, client, queries}
}

func (ds *DatabaseService) Start() error {
	// setup the database schema
	if _, err := ds.client.ExecContext(ds.Ctx, dbSchema); err != nil {
		logger.Error("error configuring db tables", "err", err)
	}

	return nil
}

func (ds *DatabaseService) Stop() error {
	return ds.client.Close()
}
