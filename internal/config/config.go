package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config defines the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

// KafkaConfig holds Kafka broker and topic configuration
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// LoadConfig reads and parses the config file
func LoadConfig() *Config {
	configFile, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer configFile.Close()

	var cfg Config
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	log.Println("Config loaded successfully!")
	return &cfg
}
