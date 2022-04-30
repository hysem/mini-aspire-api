package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config holds the configuration info
type Config struct {
	Username string
	Password string
	Host     string
	Port     uint16
	DB       string
	SSLMode  string `envconfig:"SSL_MODE"`
}

func (c Config) connStr() string {
	return fmt.Sprintf(
		"host=%[1]s port=%[2]d dbname=%[3]s user=%[4]s password=%[5]s sslmode=%[6]s",
		c.Host,
		c.Port,
		c.DB,
		c.Username,
		c.Password,
		c.SSLMode,
	)
}

func Connect(config Config) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", config.connStr())
}
