package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig *ServerConfig
	DBConfig     *DBConfig
}

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		ServerConfig: getServerConfig(),
		DBConfig:     getDBConfig(),
	}
}

func getServerConfig() *ServerConfig {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("Invalid port number")
	}
	return &ServerConfig{
		Port: port,
	}
}

func getDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnvOrDie("DB_HOST"),
		Port:     getEnvOrDie("DB_PORT"),
		Username: getEnvOrDie("DB_USER"),
		Password: getEnvOrDie("DB_PASSWORD"),
		Database: getEnvOrDie("DB_NAME"),
		SSLMode:  getEnvOrDie("DB_SSLMODE"),
	}
}

func getEnvOrDie(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Error getting environment variable %s", key)
	}
	return value
}
