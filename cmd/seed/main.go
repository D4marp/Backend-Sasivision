package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sasivision/backend/internal/config"
	"github.com/sasivision/backend/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	dir := database.MigrationsDir()
	log.Printf("Seeding database from %s ...", dir)
	if err := database.RunMigrations(db, dir); err != nil {
		log.Fatalf("seed failed: %v", err)
	}
	log.Println("Seed complete — quiz, users, markers, and videos are ready.")
}
