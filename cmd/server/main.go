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

	if cfg.RunMigrations {
		dir := database.MigrationsDir()
		log.Printf("RUN_MIGRATIONS=true — applying migrations from %s", dir)
		if err := database.RunMigrations(db, dir); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
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
