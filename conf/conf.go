package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBURI string
	Addr  string
	Debug bool
}

func New() (*Config, error) {
	var c Config

	if err := envconfig.Process("csvapp", &c); err != nil {
		return nil, fmt.Errorf("envconfig process: %w", err)
	}

	return &c, nil
}
