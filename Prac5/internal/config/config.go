package config

import (
	"os"
)

type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
	DSN      string
}

func New() Config {
	// Значения по умолчанию
	cfg := Config{
		Addr:     ":8443",
		CertFile: "certs/server.crt",
		KeyFile:  "certs/server.key",
		DSN:      "postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable",
	}

	// Переопределяем из переменных окружения
	if envAddr := os.Getenv("ADDR"); envAddr != "" {
		cfg.Addr = envAddr
	}
	if envCert := os.Getenv("CERT_FILE"); envCert != "" {
		cfg.CertFile = envCert
	}
	if envKey := os.Getenv("KEY_FILE"); envKey != "" {
		cfg.KeyFile = envKey
	}
	if envDSN := os.Getenv("DSN"); envDSN != "" {
		cfg.DSN = envDSN
	}

	return cfg
}
