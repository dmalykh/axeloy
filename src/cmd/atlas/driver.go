package atlas

import (
	"ariga.io/atlas/schema/schemaspec"
	"ariga.io/atlas/sql/migrate"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/sqlite"
)

// Driver implements the Atlas interface.
type Driver struct {
	migrate.Driver
	schemaspec.Marshaler
	schemaspec.Unmarshaler
}

func open(driverName string, db *sql.DB) (*Driver, error) {
	switch driverName {
	case "mysql":
		return mysqlProvider(db)
	case "postgres":
		return postgresProvider(db)
	case "sqlite3":
		return sqliteProvider(db)
	default:
		panic(fmt.Sprintf(`Driver "%s" doesn't supported`, driverName))
	}
}

func mysqlProvider(db *sql.DB) (*Driver, error) {
	drv, err := mysql.Open(db)
	if err != nil {
		return nil, err
	}
	return &Driver{
		Driver:      drv,
		Marshaler:   mysql.MarshalHCL,
		Unmarshaler: mysql.UnmarshalHCL,
	}, nil
}

func postgresProvider(db *sql.DB) (*Driver, error) {
	drv, err := postgres.Open(db)
	if err != nil {
		return nil, err
	}
	return &Driver{
		Driver:      drv,
		Marshaler:   postgres.MarshalHCL,
		Unmarshaler: postgres.UnmarshalHCL,
	}, nil
}

func sqliteProvider(db *sql.DB) (*Driver, error) {
	drv, err := sqlite.Open(db)
	if err != nil {
		return nil, err
	}
	return &Driver{
		Driver:      drv,
		Marshaler:   sqlite.MarshalHCL,
		Unmarshaler: sqlite.UnmarshalHCL,
	}, nil
}
