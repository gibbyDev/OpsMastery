package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gibbyDev/OpsMastery/handlers"
    "github.com/gibbyDev/OpsMastery/middleware"
    "gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    api := app.Group("/api/v1")
    {
        api.Get("/", func(c *fiber.Ctx) error {
            return c.SendString("Server is running!")
        })

        api.Post("/signup", handlers.SignUp)
        api.Post("/signin", handlers.SignIn)

        api.Use(middleware.JWTMiddleware)

        api.Post("/signout", handlers.SignOut)

        api.Delete("/users/:id", middleware.OnlyAdmin(db, handlers.DeleteUserByID))
        api.Put("/users/:id/role", middleware.OnlyAdmin(db, handlers.SetUserRole))

        api.Get("/users", middleware.OnlyModerator(db, handlers.ListUsers))
        api.Get("/users/:id", middleware.OnlyModerator(db, handlers.GetUserByID))
        api.Put("/users/:id", middleware.OnlyModerator(db, handlers.UpdateUserByID))

        api.Get("/users/me", middleware.OnlyUser(db, handlers.GetCurrentUser))
        api.Put("/users/me", middleware.OnlyUser(db, handlers.UpdateCurrentUser)) 

        api.Post("/ticket", handlers.CreateTicket)
        api.Get("/tickets", handlers.ListTickets)
        api.Get("/ticket/:id", handlers.GetTicketByID)
        api.Put("/ticket/:id", handlers.UpdateTicketByID)
        api.Delete("/ticket/:id", handlers.DeleteTicketByID)
    }
}