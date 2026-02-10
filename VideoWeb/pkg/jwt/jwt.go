package jwt

import "github.com/golang-jwt/jwt"

var secretKey = []byte("your-secret-key")

type Claims struct {
	UserID uint `json:"user-id"`
	jwt.StandardClaims
}
