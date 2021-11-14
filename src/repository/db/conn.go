package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	"log"
	"os"
	"strings"
)

func Connect(ctx context.Context, drivenName string, dsn string) (*reform.DB, error) {
	sqlDB, err := sql.Open(drivenName, dsn)
	if err != nil {
		return nil, fmt.Errorf(`(sql open "%s") %w`, drivenName, err)
	}
	go func() {
		<-ctx.Done()
		sqlDB.Close()
	}()

	// Use new *log.Logger for logging.
	logger := log.New(os.Stderr, "SQL: ", log.Flags())

	return reform.NewDB(sqlDB, DialectFor(drivenName), reform.NewPrintfLogger(logger.Printf)), nil

}

// DialectFor returns reform Dialect for given driver string, or nil.
func DialectFor(driver string) reform.Dialect {
	// for sqlite3_with_sleep
	if strings.HasPrefix(driver, "sqlite3") {
		return sqlite3.Dialect
	}

	switch driver {
	case "postgres", "pgx":
		return postgresql.Dialect
	case "mysql":
		return mysql.Dialect
	case "sqlserver":
		return sqlserver.Dialect
	default:
		return nil
	}
}
