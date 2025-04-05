package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    "OpsMastery/user-service/initialization"
    "OpsMastery/user-service/routes"
    "OpsMastery/user-service/handlers"
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

    routes.SetupUserRoutes(app, db)

    log.Println("User Service running on :8082")
    log.Fatal(app.Listen(":8082"))
}