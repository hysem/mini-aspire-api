package config

import (
	"github.com/hysem/mini-aspire-api/app/core/bcrypt"
	"github.com/hysem/mini-aspire-api/app/core/db"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
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
	Database struct {
		Master db.Config
	}
	Bcrypt bcrypt.Config
	JWT    jwt.Config
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
