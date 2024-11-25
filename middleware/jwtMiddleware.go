package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gibbyDev/OpsMastery/utils"
    "log"
)

func JWTMiddleware(c *fiber.Ctx) error {
    token := c.Cookies("jwt")

    if token == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }

    claims, err := utils.ValidateJWT(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }

    log.Printf("JWT Claims: %+v\n", claims)

    c.Locals("userID", claims["sub"])   
    c.Locals("userRole", claims["role"]) 

    return c.Next() 
}