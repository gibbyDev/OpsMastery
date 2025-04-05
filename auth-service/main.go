package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"OpsMastery/auth-service/initialization"
	"OpsMastery/auth-service/routes"
	"OpsMastery/auth-service/handlers"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	initialization.LoadEnv()
	db := initialization.SetupDatabase()

	    // Pass the database instance to the handlers package
    handlers.SetDB(db)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	routes.SetupAuthRoutes(app, db)

	log.Println("Auth Service running on :8081")
	log.Fatal(app.Listen(":8081"))
}