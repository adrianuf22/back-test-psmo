package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host         string `required:"true"`
	Port         int    `default:"5432"`
	User         string `required:"true"`
	Password     string `required:"true"`
	DB           string `required:"true"`
	Driver       string `default:"postgres"`
	SSLMode      string `split_words:"true" default:"disable"`
	MaxOpenConns int    `default:"5"`
	MaxIdleConns int    `default:"5"`
	Timeout      int    `default:"15000"`
}

func NewDatabase() Database {
	var db Database
	envconfig.MustProcess("POSTGRES", &db)

	return db
}

func (d *Database) DSN() string {
	return fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s&user=%s&password=%s",
		d.Host,
		d.Port,
		d.DB,
		d.SSLMode,
		d.User,
		d.Password,
	)
}
