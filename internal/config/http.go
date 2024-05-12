package config

import "github.com/kelseyhightower/envconfig"

type Http struct {
	MaxBytesReader    int    `default:"1048576" split_words:"true"`
	ReadHeaderTimeout int    `default:"10000" split_words:"true"`
	Port              string `default:"3000"`
}

func NewHttp() Http {
	var h Http
	envconfig.MustProcess("HTTP", &h)

	return h
}
