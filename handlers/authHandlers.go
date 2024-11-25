package handlers

import (
    "log"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gibbyDev/OpsMastery/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/gibbyDev/OpsMastery/utils"
	"time"
)

func SignUp(c *fiber.Ctx) error {
    var user models.User
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Name     string `json:"name"`
        Role     string `json:"role"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    user.Email = input.Email
    user.Password = input.Password
    user.Name = input.Name
    user.Role = input.Role

    storedHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }
    user.Password = string(storedHash)

    log.Printf("Hashed password for user %s: %s\n", user.Email, user.Password)

    if err := db.Create(&user).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(http.StatusCreated).JSON(user)
}

func SignIn(c *fiber.Ctx) error {
    var userInput struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var user models.User
    
    if err := c.BodyParser(&userInput); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := db.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
        log.Println("User not found:", err)
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
    }

    c.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    token,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true, 
    })
    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Sign in successful"})
}

func SignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), 
		HTTPOnly: true,
	})
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully signed out"})
}
