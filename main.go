package main

import (
    "log"
    "OpsMastery/initialization"
    "OpsMastery/handlers"
    "OpsMastery/routes"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

var db *gorm.DB

func init() {
    initialization.LoadEnv()
    db = initialization.SetupDatabase()
    if db == nil {
        log.Fatal("Database connection is nil after initialization")
    }
    handlers.SetDB(db)
}

func main() {
    app := fiber.New()

    routes.SetupRoutes(app, db)

    log.Println("Starting server on :8080")
    if err := app.Listen(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}