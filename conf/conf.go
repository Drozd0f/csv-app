package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBURI     string
	Addr      string `default:":8080"`
	Debug     bool
	ChunkSize int32 `split_words:"true" default:"100"`
}

func New() (*Config, error) {
	var c Config

	if err := envconfig.Process("csvapp", &c); err != nil {
		return nil, fmt.Errorf("envconfig process: %w", err)
	}

	return &c, nil
}
