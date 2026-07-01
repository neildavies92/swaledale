package config

import "os"

type Config struct {
	Addr          string
	DatabaseURL   string
	SecureCookies bool
}

func Load() Config {
	return Config{
		Addr:          getEnv("API_ADDR", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://swaledale:swaledale@localhost:5432/swaledale?sslmode=disable"),
		SecureCookies: getEnv("SECURE_COOKIES", "false") == "true",
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
