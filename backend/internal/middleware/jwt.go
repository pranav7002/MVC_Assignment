package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error": "Invalid token claims"}`, http.StatusUnauthorized)
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, `{"error": "Invalid user_id claim"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("user_id"), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
