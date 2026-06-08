package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		return Config{
			ConnString: "postgres://mirudull:1234567890@192.168.0.108:5432/cloudhub?sslmode=disable",
			Port:       ":8080",
		}
	}
	return Config{
		ConnString: getEnv("connString",
			"postgres://mirudull:1234567890@192.168.0.108:5432/cloudhub?sslmode=disable"),
		Port: getEnv("port", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

type Config struct {
	ConnString string
	Port       string
}
