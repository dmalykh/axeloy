package db

import (
	"context"
	"database/sql"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"log"
	"os"
)

func Connect(ctx context.Context, drivenName string, dsn string) (*reform.DB, error) {
	sqlDB, err := sql.Open(drivenName, dsn)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		sqlDB.Close()
	}()

	// Use new *log.Logger for logging.
	logger := log.New(os.Stderr, "SQL: ", log.Flags())

	return reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf)), nil

}
