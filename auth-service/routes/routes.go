package routes

import (
    "github.com/gofiber/fiber/v2"
    "OpsMastery/auth-service/handlers"
    "gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {

    handlers.SetDB(db) // Pass the database instance to the handlers package

    api := app.Group("/api/v1")

    // Public routes (no authentication required)
    api.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Auth Service is running!")
    })
    api.Post("/signup", handlers.SignUp)
    api.Post("/signin", handlers.SignIn)
    api.Get("/verify/:token", handlers.VerifyEmail)
    api.Post("/forgot-password", handlers.RequestPasswordReset)
    api.Post("/reset-password", handlers.ResetPassword)

    // Protected routes (require authentication)
    protected := api.Group("")
    protected.Post("/signout", handlers.SignOut)
    protected.Post("/auth/refresh", handlers.RefreshToken)
}