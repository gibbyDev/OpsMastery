package routes

import (
    "github.com/gofiber/fiber/v2"
    "OpsMastery/user-service/handlers"
    "OpsMastery/user-service/middleware"
    "gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
    api := app.Group("/api/v1")

    // Protected routes (require authentication)
    protected := api.Group("")
    protected.Use(middleware.JWTMiddleware)

    // Admin routes
    protected.Delete("/users/:id", middleware.OnlyAdmin(db, handlers.DeleteUserByID))
    protected.Put("/users/:id/role", middleware.OnlyAdmin(db, handlers.SetUserRole))

    // Moderator routes
    protected.Get("/users", middleware.OnlyModerator(db, handlers.ListUsers))
    protected.Get("/users/:id", middleware.OnlyModerator(db, handlers.GetUserByID))
    protected.Put("/users/:id", middleware.OnlyModerator(db, handlers.UpdateUserByID))

    // User routes
    protected.Get("/users/me", middleware.OnlyUser(db, handlers.GetCurrentUser))
    protected.Put("/users/me", middleware.OnlyUser(db, handlers.UpdateCurrentUser))
}