package middleware

import (
	"context"
	"net/http"
	jj "newww/internal/middleware/jwt"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := jj.ParseAccessToken(cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), jj.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
