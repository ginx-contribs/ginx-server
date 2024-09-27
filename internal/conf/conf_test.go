package conf

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestReadFrom(t *testing.T) {
	filename := "testdata/conf.toml"
	app, err := ReadFrom(filename)
	assert.NoError(t, err)
	t.Log(app)
}

func TestRevise(t *testing.T) {
	cfg := App{Server: Server{Address: "127.0.0.1:8080"}, Log: Log{Level: slog.LevelDebug}}
	reviseConf, err := Revise(cfg)
	assert.NoError(t, err)
	t.Log(reviseConf)
}

func TestWriteTo(t *testing.T) {
	filename := "testdata/conf.toml"
	err := WriteTo(filename, DefaultConfig)
	assert.NoError(t, err)
}
