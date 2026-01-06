package main

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

// Config maps the structure of config.yaml
type Config struct {
	Server struct {
		Port          int    `yaml:"port"`
		Host          string `yaml:"host"`
		Tls_Cert_Path string `yaml:"tls-cert-path"`
		Tls_Key_Path  string `yaml:"tls-key-path"`
	} `yaml:"server"`

	Database struct {
		Port     int    `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	AllowedOrigins []string `yaml:"allowed_origins"`
}

func GetOptions() *Config {

	// Open the YAML file
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create a variable to store decoded data
	var cfg Config

	// Unmarshal YAML into the struct
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	return &cfg
}
