package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpire  string
	RefreshExpire string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type RabbitMQConfig struct {
	URL           string
	TransferQueue string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "ewallet_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "25c1f72ff7ad9fae385328d34cce961452cd0209e3992c0c64866929be371c45c0a95384f81ec7e74cb36e780de7fbde99ba266a12ba640e82a39468ef8ae27a6ada4bf68b706d639a8b2dc2858177fc75f6c20a89c3aa416373f9f1f67ecf18237abbf4b9e5af5946fa4ed2b982a0765335a991e70557be8a3e343b110b1234314b152767918970bc2fb5967488345777a76dc2aac546c33dc15db1438eeade5e96b2aa5b6411cc8c0fb51c85e12e811afbda54c1d0fabac3d742cb04829b10928088f20c224ebd564a82a9830e4dd333967cf6147a54198dbdc267ce2f295fbb87dd1cc5d50c4fa8ee850fabf775925ed84f65c700e81e383f14ae36582c35"),
			AccessExpire:  getEnv("JWT_ACCESS_EXPIRE", "24h"),
			RefreshExpire: getEnv("JWT_REFRESH_EXPIRE", "168h"),
		},
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:           getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			TransferQueue: getEnv("TRANSFER_QUEUE", "transfer_queue"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
