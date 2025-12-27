package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID int `json:"user_id"`
	Admin  int `json:"admin"`
	jwt.RegisteredClaims
}
