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
    user.Active = false
    user.VerificationToken = utils.GenerateRandomToken()

    storedHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }
    user.Password = string(storedHash)

    log.Printf("Hashed password for user %s: %s\n", user.Email, user.Password)

    if err := db.Create(&user).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    if err := utils.SendVerificationEmail(user.Email, user.VerificationToken); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send verification email"})
    }

    return c.Status(http.StatusCreated).JSON(fiber.Map{
        "message": "Registration successful. Please check your email to verify your account.",
    })
}

func SignIn(c *fiber.Ctx) error {
    var userInput struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&userInput); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    var user models.User
    if err := db.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    if !user.Active {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Please verify your email before signing in",
        })
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    accessToken, refreshToken, err := utils.GenerateJWT(user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate tokens"})
    }

    log.Printf("Setting access token: %s", accessToken)
    log.Printf("Setting refresh token: %s", refreshToken)

    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Expires:  time.Now().Add(15 * time.Minute),
        HTTPOnly: true,
        Secure:   false,
        SameSite: "None",
        Domain:   "",
        Path:     "/",
    })

    c.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Expires:  time.Now().Add(7 * 24 * time.Hour),
        HTTPOnly: true,
        Secure:   false,
        SameSite: "None",
        Domain:   "",
        Path:     "/",
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

func RefreshToken(c *fiber.Ctx) error {
    refreshToken := c.Cookies("refresh_token")
    if refreshToken == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Refresh token not found",
        })
    }

    claims, err := utils.ValidateJWT(refreshToken, true)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid refresh token",
        })
    }

    // Get user from claims
    var user models.User
    if err := db.First(&user, claims["sub"]).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    // Generate new tokens
    accessToken, refreshToken, err := utils.GenerateJWT(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not generate tokens",
        })
    }

    // Set new cookies
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Expires:  time.Now().Add(15 * time.Minute),
        HTTPOnly: true,
    })

    c.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Expires:  time.Now().Add(7 * 24 * time.Hour),
        HTTPOnly: true,
    })

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Tokens refreshed successfully",
    })
}

func VerifyEmail(c *fiber.Ctx) error {
    token := c.Params("token")

    var user models.User
    if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Invalid verification token",
        })
    }

    user.Active = true
    user.VerificationToken = "" // Clear the token after verification

    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to verify email",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Email verified successfully",
    })
}

func RequestPasswordReset(c *fiber.Ctx) error {
    var input struct {
        Email string `json:"email"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    var user models.User
    if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        // Don't reveal if email exists or not
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "message": "If your email is registered, you will receive a password reset link",
        })
    }

    resetToken := utils.GenerateRandomToken()
    user.ResetToken = resetToken
    user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)

    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to process password reset",
        })
    }

    if err := utils.SendPasswordResetEmail(user.Email, resetToken); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to send reset email",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "If your email is registered, you will receive a password reset link",
    })
}

func ResetPassword(c *fiber.Ctx) error {
    token := c.Params("token")
    var input struct {
        Password string `json:"password"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    var user models.User
    if err := db.Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid or expired reset token",
        })
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to process new password",
        })
    }

    user.Password = string(hashedPassword)
    user.ResetToken = ""
    user.ResetTokenExpiry = time.Time{}

    if err := db.Save(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update password",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Password has been reset successfully",
    })
}

