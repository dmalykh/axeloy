package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/xo/dburl"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	"log"
	"os"
	"strings"
)

// Connect to database with driver and dsn
func Connect(ctx context.Context, drivenName string, dsn string) (*sql.DB, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf(`wrong dsn %w`, err)
	}

	db, err := sql.Open(drivenName, u.DSN)
	if err != nil {
		return nil, fmt.Errorf(`(sql open "%s") %w`, drivenName, err)
	}
	//Shutdown database connection
	go func() {
		<-ctx.Done()
		db.Close()
	}()
	return db, nil
}

//
func Reform(driverName string, db *sql.DB) *reform.DB {
	// Use new *log.Logger for logging.
	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	return reform.NewDB(db, DialectFor(driverName), reform.NewPrintfLogger(logger.Printf))

}

// DialectFor returns reform Dialect for given driver string, or nil.
func DialectFor(driverName string) reform.Dialect {
	if strings.HasPrefix(driverName, "sqlite3") {
		return sqlite3.Dialect
	}

	switch driverName {
	case "postgres":
		return postgresql.Dialect
	case "mysql":
		return mysql.Dialect
	case "sqlserver":
		return sqlserver.Dialect
	default:
		panic(fmt.Sprintf(`Driver "%s" doesn't supported`, driverName))
	}
}
