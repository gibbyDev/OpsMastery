package handlers

import (
    "net/http"
    "strconv"
    "github.com/gofiber/fiber/v2"
    "github.com/gibbyDev/OpsMastery/models"
)

func ListUsers(c *fiber.Ctx) error {
    var users []models.User
    if err := db.Unscoped().Find(&users).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
    id := c.Params("id")

    var user models.User
    if err := db.First(&user, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }
    return c.JSON(user)
}

func DeleteUserByID(c *fiber.Ctx) error {
    id := c.Params("id")

    userIDParsed, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    result := db.Unscoped().Delete(&models.User{}, userIDParsed)
    if result.RowsAffected == 0 {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }
    if result.Error != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
    }
    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

func UpdateUserByID(c *fiber.Ctx) error {
    id := c.Params("id")

    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    if err := db.Model(&user).Where("id = ?", id).Updates(user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
    }
    return c.JSON(user)
}

func SetUserRole(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User

    if err := db.First(&user, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    var input struct {
        Role string `json:"role"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    user.Role = input.Role
    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user role"})
    }
    return c.JSON(user)
}

func GetCurrentUser(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(uint)
    
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
    }

    var user models.User
    if err := db.First(&user, userID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }
    return c.JSON(user)
}

func UpdateCurrentUser(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
    }

    var user models.User
    if err := db.First(&user, userID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
    }
    return c.JSON(user)
}
