package auth

import (
	"net/http"
	"newww/internal/middleware/jwt"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "no refresh token", http.StatusUnauthorized)
		return
	}

	claims, err := jwt.ParseRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Генерация нового access token
	accessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.Admin)

	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	// Ставим новый access token в куки
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
