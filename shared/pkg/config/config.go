package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	HEREAPIKey string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int // hours
}

func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	// Check if DATABASE_URL is provided (for container environments)
	databaseURL := getEnv("DATABASE_URL", "")
	var dbConfig DatabaseConfig

	if databaseURL != "" {
		// Parse DATABASE_URL if provided
		dbConfig = parseDatabaseURL(databaseURL)
	} else {
		// Use individual environment variables
		dbConfig = DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "video_streaming"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}
	}

	return &Config{
		Database: dbConfig,
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		HEREAPIKey: getEnv("HERE_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// parseDatabaseURL parses a DATABASE_URL string into DatabaseConfig
// Format: postgres://user:password@host:port/dbname?sslmode=disable
func parseDatabaseURL(databaseURL string) DatabaseConfig {
	u, err := url.Parse(databaseURL)
	if err != nil {
		// Return default config if parsing fails
		return DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "video_streaming",
			SSLMode:  "disable",
		}
	}

	// Extract components
	host := u.Hostname()
	port := 5432
	if u.Port() != "" {
		if p, err := strconv.Atoi(u.Port()); err == nil {
			port = p
		}
	}

	user := "postgres"
	password := "password"
	if u.User != nil {
		user = u.User.Username()
		if pwd, ok := u.User.Password(); ok {
			password = pwd
		}
	}

	dbname := strings.TrimPrefix(u.Path, "/")
	if dbname == "" {
		dbname = "video_streaming"
	}

	sslmode := "disable"
	if u.RawQuery != "" {
		params, _ := url.ParseQuery(u.RawQuery)
		if ssl := params.Get("sslmode"); ssl != "" {
			sslmode = ssl
		}
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}
}
