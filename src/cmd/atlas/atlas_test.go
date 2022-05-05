package atlas

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAtlas_Inspect(t *testing.T) {
	dir, _ := os.Getwd()
	var dbpath = dir + "/sqlite-test-database.db"
	var dsn = createdb(t, dbpath)
	defer os.Remove(dbpath)

	// Fill databse
	func() {
		db, err := sql.Open("sqlite3", dbpath)
		assert.NoError(t, err)
		_, err = db.Exec(`CREATE TABLE student ("idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,"name" TEXT);`)
		assert.NoError(t, err)
		db.Close()
	}()

	var ctx = context.TODO()
	atlas, err := NewAtlas(ctx, `sqlite3`, dsn)
	assert.NoError(t, err)

	data, err := atlas.Inspect(ctx)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "idStudent")
}

func TestAtlas_ApplyNew(t *testing.T) {
	dir, _ := os.Getwd()
	var dbpath = dir + "/sqlite-test-database.db"
	var dsn = createdb(t, dbpath)
	defer os.Remove(dbpath)

	var ctx = context.TODO()
	atlas, err := NewAtlas(ctx, `sqlite3`, dsn)
	assert.NoError(t, err)

	// Create migration file
	var migfile = dir + "/test-migration.hcl"
	defer os.Remove(migfile)
	func() {
		assert.NoError(t, os.WriteFile(migfile, []byte(
			`table "categories" {
			  schema = schema.main
			  column "name" {
				null = true
				type = text
			  }
			  column "ohmysupername" {
				null           = false
				type           = integer
				auto_increment = true
			  }
			}
			schema "main" {
			}`), 777))
	}()
	ok, err := atlas.Apply(ctx, migfile)
	assert.NoError(t, err)
	assert.True(t, ok)

	data, err := atlas.Inspect(ctx)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "ohmysupername")
}

// Create SQLite file, fill it and return dsn
func createdb(t *testing.T, path string) string {
	os.Remove(path)
	file, err := os.Create(path)
	defer file.Close()
	assert.NoError(t, err)

	var dsn = fmt.Sprintf(`sqlite://file:%s`, path)
	return dsn
}
