package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateToken(username string, email string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": username,                           // Subject (user identifier)
		"eml": email,                              // User email
		"iss": "golang-microservice-sekolah",      // Issuer
		"exp": time.Now().AddDate(1, 0, 0).Unix(), // Expiration time
		"iat": time.Now().Unix(),                  // Issued at
	})
	
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	
	// Print information about the created token
	log.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}
