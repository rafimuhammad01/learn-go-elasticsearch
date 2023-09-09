package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CloudID string
	APIKey  string
}

func NewConfig(file string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("env file doesn't exists: %s", file)
	}

	return &Config{
		CloudID: os.Getenv("CLOUD_ID"),
		APIKey:  os.Getenv("API_KEY"),
	}, nil
}
