package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config holds the configuration info
type Config struct {
	App struct {
		Port        int
		Environment string
		Docs        string
	}
	Log struct {
		Level string
	}
}

var (
	// current holds the currently loaded configuration
	current Config
)

// Load parses and load the configuration this must be called first
func Load() error {
	if err := envconfig.Process("MINI_ASPIRE_API", &current); err != nil {
		return errors.Wrap(err, "failed to parse environment variables")
	}
	return nil
}

// Current returns the currently loaded configuration
func Current() Config {
	return current
}
