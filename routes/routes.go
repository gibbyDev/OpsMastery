package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gibbyDev/OpsMastery/handlers"
    "github.com/gibbyDev/OpsMastery/middleware"
    "gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    api := app.Group("/api/v1")
    
    // Public routes (no authentication required)
    api.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Server is running!")
    })
    api.Post("/signup", handlers.SignUp)
    api.Post("/signin", handlers.SignIn)
    api.Get("/verify/:token", handlers.VerifyEmail)
    api.Post("/forgot-password", handlers.RequestPasswordReset)
    api.Post("/reset-password/:token", handlers.ResetPassword)

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

    // Other protected routes
    protected.Post("/signout", handlers.SignOut)
    protected.Post("/auth/refresh", handlers.RefreshToken)
    protected.Post("/ticket", handlers.CreateTicket)
    protected.Get("/tickets", handlers.ListTickets)
    protected.Get("/ticket/:id", handlers.GetTicketByID)
    protected.Put("/ticket/:id", handlers.UpdateTicketByID)
    protected.Delete("/ticket/:id", handlers.DeleteTicketByID)
}