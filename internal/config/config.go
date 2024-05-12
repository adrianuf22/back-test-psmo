package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	Http
	Database
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func New() *Config {
	err := godotenv.Load(filepath.Join(basepath, "../../.env"))
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Http:     NewHttp(),
		Database: NewDatabase(),
	}
}
