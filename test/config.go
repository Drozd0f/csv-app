package test

import (
	"testing"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/test/containers"
)

func NewConfig(t *testing.T, td *containers.TestDatabase) *conf.Config {
	return &conf.Config{
		DBURI:     td.ConnectionString(t),
		Addr:      "",
		Debug:     false,
		ChunkSize: 10,
	}
}
