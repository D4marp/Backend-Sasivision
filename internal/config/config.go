package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	// Server
	AppEnv  string
	AppPort string
	AppName string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret string
	JWTExpiry time.Duration

	// CORS
	CORSOrigin string

	// Rate limiting
	RateLimitRPS int
}

func LoadConfig() *Config {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "720h"))

	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))

	return &Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		AppPort:      getEnv("APP_PORT", "8080"),
		AppName:      getEnv("APP_NAME", "SasiVision-API"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "sasivision"),
		JWTSecret:    getEnv("JWT_SECRET", "your_secret_key"),
		JWTExpiry:    jwtExpiry,
		CORSOrigin:   getEnv("CORS_ORIGIN", "http://localhost:8000"),
		RateLimitRPS: rateLimitRPS,
	}
}

func InitDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
