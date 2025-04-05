package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"OpsMastery/ticket-service/initialization"
	"OpsMastery/ticket-service/routes"
	"OpsMastery/ticket-service/handlers"
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

	routes.SetupTicketRoutes(app, db)

	log.Println("Ticket Service running on :8083")
	log.Fatal(app.Listen(":8083"))
}