package main

import (
    "log"
    "github.com/gibbyDev/OpsMastery/initialization"
    "github.com/gibbyDev/OpsMastery/handlers"
    "github.com/gibbyDev/OpsMastery/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "gorm.io/gorm"
    "os"
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

    app.Use(cors.New(cors.Config{
        AllowOrigins:     os.Getenv("FRONTEND_URL"),
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
        MaxAge:           300,
    }))

    routes.SetupRoutes(app, db)

    log.Println("Starting server on :8080")
    if err := app.Listen(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
