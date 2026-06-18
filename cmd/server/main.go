package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sasivision/backend/internal/config"
	"github.com/sasivision/backend/internal/database"
	"github.com/sasivision/backend/internal/handlers"
	"github.com/sasivision/backend/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.EnsureSchema(db, cfg.RunMigrations); err != nil {
		log.Fatalf("Database bootstrap failed: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	h := handlers.New(db, cfg)
	router := server.SetupRouter(h, cfg)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting %s API server on %s", cfg.AppName, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
