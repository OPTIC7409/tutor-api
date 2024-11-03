package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}, nil
}

func loadEnv() error {
	err := godotenv.Load()
	if err == nil {
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	rootDir, err := findRootDir(currentDir)
	if err != nil {
		return err
	}

	return godotenv.Load(filepath.Join(rootDir, ".env"))
}

func findRootDir(currentDir string) (string, error) {
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", os.ErrNotExist
		}
		currentDir = parentDir
	}
}
