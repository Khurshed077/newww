package middleware

import (
	"context"
	"net/http"
	"strings"

	"newww/internal/middleware/jwt"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tokenString string

		// 1️⃣ Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 2️⃣ fallback: cookie access_token
			cookie, err := r.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			}
		}

		// 3️⃣ ЕСЛИ токена нет — это ОК (публичный запрос)
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 4️⃣ Парсим access token
		claims, err := jwt.ParseAccessToken(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 5️⃣ Кладём claims в context
		ctx := context.WithValue(r.Context(), jwt.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
