package main

import (
	"flag"
	"fmt"
	"log"

	"meal-planner/internal/auth"
	"meal-planner/internal/config"
	"meal-planner/internal/db"
	"meal-planner/internal/models"
)

func main() {
	name := flag.String("name", "", "Display name for the new user (required)")
	isAdmin := flag.Bool("admin", false, "Grant the new user admin privileges")
	flag.Parse()

	if *name == "" {
		log.Fatal("Usage: go run cmd/createuser/main.go -name \"Full Name\" [-admin]")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	repo := db.NewRepository(database)

	apiKey, err := auth.GenerateAPIKey()
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}

	deviceID, err := auth.GenerateDeviceID()
	if err != nil {
		log.Fatalf("Failed to generate device ID: %v", err)
	}

	user := &models.User{
		Name:     *name,
		DeviceID: deviceID,
		APIKey:   apiKey,
		IsAdmin:  *isAdmin,
	}

	user, err = repo.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("User created: %s (admin: %t)\n", user.Name, user.IsAdmin)
	fmt.Printf("API key (save it now, it won't be shown again): %s\n", user.APIKey)
}
