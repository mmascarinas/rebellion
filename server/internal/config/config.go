package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        		 string
	Environment 		 string
	AllowedOrigins   string
	RedisAddr        string
	RedisPassword    string
	RedisCACert      string
	AuthClientID     string
	AuthClientSecret string
	AuthRedirectURL  string
	AuthUserScopeURL string
	AuthUserOauthURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:           	getEnv("PORT", "8080"),
		Environment:    	getEnv("ENVIRONMENT", "development"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "*"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisCACert:     getEnv("REDIS_CA_CERT", ""),
		AuthClientID:     getEnv("AUTH_CLIENT_ID", ""),
		AuthClientSecret: getEnv("AUTH_CLIENT_SECRET", ""),
		AuthRedirectURL:  getEnv("AUTH_REDIRECT_URL", "http://localhost:8080/auth/login"),
		AuthUserScopeURL: getEnv("AUTH_USER_SCOPE_URL", ""),
		AuthUserOauthURL: getEnv("AUTH_USER_OAUTH_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
