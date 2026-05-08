package router

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"

	"rebellion/internal/config"
	"rebellion/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Setup(cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	redisOpts := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	}

	if cfg.Environment == "production" && cfg.RedisCACert != "" {
		caCert, err := os.ReadFile(cfg.RedisCACert)
		if err != nil {
			log.Fatalf("failed to read Redis CA cert: %v", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		redisOpts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		}
	}

	rdb := redis.NewClient(redisOpts)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     parseOrigins(cfg.AllowedOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	auth := handler.NewAuthHandler(rdb, cfg)
	gameH := handler.NewGameHandler()

	r.GET("/health", handler.HealthCheck)
	r.GET("/auth/login", auth.Login)
	r.POST("/auth/logout", auth.Logout)
	r.GET("/ws/game/create", gameH.CreateGame)
	r.GET("/ws/game/:id", gameH.JoinGame)

	return r
}

func parseOrigins(origins string) []string {
	parts := strings.Split(origins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
