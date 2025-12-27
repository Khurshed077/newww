package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("SECRET_KEY")
var refreshSecret = []byte("REFRESH_SECRET_KEY") // отдельный секрет для refresh

func GenerateAccessToken(userID int, admin int) (string, error) {
	claims := jwtlib.MapClaims{
		"user_id": userID,
		"admin":   admin,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // короткий срок
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(userID int, admin int) (string, error) {
	claims := jwtlib.MapClaims{
		"user_id": userID,
		"admin":   admin,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}
func ParseAccessToken(tokenStr string) (*Claims, error) {
	claims := jwtlib.MapClaims{}
	token, err := jwtlib.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	userID := int(claims["user_id"].(float64))
	admin := 0
	if a, ok := claims["admin"]; ok {
		if af, ok := a.(float64); ok {
			admin = int(af)
		}
	}

	return &Claims{
		UserID: userID,
		Admin:  admin,
	}, nil
}

func ParseRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwtlib.Parse(tokenStr, func(t *jwtlib.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	mc, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID := int(mc["user_id"].(float64))
	admin := 0
	if a, ok := mc["admin"]; ok {
		if af, ok := a.(float64); ok {
			admin = int(af)
		}
	}

	return &Claims{
		UserID: userID,
		Admin:  admin,
	}, nil
}
