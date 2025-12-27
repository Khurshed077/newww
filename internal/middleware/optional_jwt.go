package middleware

import (
	"context"
	"net/http"
	"newww/internal/middleware/jwt"
)

func OptionalJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			// гость
			next.ServeHTTP(w, r)
			return
		}

		claims, err := jwt.ParseAccessToken(cookie.Value)
		if err != nil {
			// битый токен
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
