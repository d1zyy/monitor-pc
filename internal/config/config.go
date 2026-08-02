package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
	Addr string
}

func Load() (Config, error) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "0.0.0.0"
	}

	if port == "" {
		port = "8080"
	}

	return Config{
		Host: host,
		Port: port,
		Addr: host + ":" + port,
	}, nil
}
