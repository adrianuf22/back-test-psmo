package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Dsn          string `required:"true"`
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
