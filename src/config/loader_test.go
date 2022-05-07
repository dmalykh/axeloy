package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLoad(t *testing.T) {

	var config = `database "pgx" {
		  dsn = "postgre://name@pass:serv"
		}
		
		driver "graphql" "ways/graphql/graphql.so" {
		  listen_addr = "127.0.0.1:8080"
		}
		
		driver "superdemo" "ways/graphql/graphql.so" {
		  port = "998"
		}`
	var filepath = t.TempDir() + "/test.hcl"
	assert.NoError(t, ioutil.WriteFile(filepath, []byte(config), 777))

	c, err := LoadFile(filepath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(c.Drivers))
	assert.Equal(t, c.Database.Driver, "pgx")

	// Check body decoding
	t.Run(`Check body decoding`, func(t *testing.T) {
		type Zero struct {
			ListenAddr string `hcl:"listen_addr"`
		}
		var k Zero
		err := Unmarshal(c.Drivers[0], &k)
		assert.NoError(t, err)
		assert.Equal(t, `127.0.0.1:8080`, k.ListenAddr)
	})
}
