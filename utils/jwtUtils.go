package utils

import (
	"log"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gibbyDev/OpsMastery/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,                         
		"email": user.Email,                      
		"role":  user.Role,                       
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorClaimsInvalid)
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
