package config_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veleton777/booking_api/internal/config"
)

func TestConfig(t *testing.T) {
	env := map[string]string{
		"APP_NAME":  "app",
		"LOG_LEVEL": "3",

		"HTTP_PORT": "9090",
	}

	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			require.NoError(t, err)
		}
	}

	conf, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, conf.AppName(), "app")
	assert.Equal(t, conf.LogLevel(), zerolog.Level(3))
	assert.Equal(t, conf.HTTPAddr(), ":9090")
}
