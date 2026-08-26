package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
	Addr string

	PprofEnabled bool
	PprofPort    string
	PprofHost    string
	PprofAddr    string
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

	pprofEnabled := os.Getenv("PPROF_ENABLED")
	pprofPort := os.Getenv("PPROF_PORT")
	pprofHost := os.Getenv("PPROF_HOST")

	if pprofPort == "" {
		pprofPort = "6060"
	}

	if pprofHost == "" {
		pprofHost = "127.0.0.1"
	}

	isProfEnabled := pprofEnabled == "true"

	return Config{
		Host: host,
		Port: port,
		Addr: host + ":" + port,

		PprofEnabled: isProfEnabled,
		PprofPort:    pprofPort,
		PprofHost:    pprofHost,
		PprofAddr:    pprofHost + ":" + pprofPort,
	}, nil
}
