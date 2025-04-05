package routes

import (
    "github.com/gofiber/fiber/v2"
    "OpsMastery/ticket-service/handlers"
    "OpsMastery/ticket-service/middleware"
    "gorm.io/gorm"
)

func SetupTicketRoutes(app *fiber.App, db *gorm.DB) {
    api := app.Group("/api/v1")

    // Protected routes (require authentication)
    protected := api.Group("")
    protected.Use(middleware.JWTMiddleware)

    // Ticket routes
    protected.Post("/ticket", handlers.CreateTicket)
    protected.Get("/tickets", handlers.ListTickets)
    protected.Get("/ticket/:id", handlers.GetTicketByID)
    protected.Put("/ticket/:id", handlers.UpdateTicketByID)
    protected.Delete("/ticket/:id", handlers.DeleteTicketByID)
}