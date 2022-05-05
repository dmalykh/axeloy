package app

import (
	configuration "github.com/dmalykh/axeloy/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_Load(t *testing.T) {
	var a = new(App)
	var ctx = a.WithShutdown()

	type testCase struct {
		name   string
		config *configuration.Config
	}
	var testCases = []testCase{}
	for _, tc := range testCases {
		_, err := a.Load(ctx, tc.config)
		assert.NoError(t, err)
	}
}
