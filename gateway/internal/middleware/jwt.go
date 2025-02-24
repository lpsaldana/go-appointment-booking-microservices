package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secretKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Auth token needed", http.StatusUnauthorized)
			return
		}

		// Header must be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// sign HS256 validation
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected jwt sign method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			log.Printf("Error validatin token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Valid token
		next.ServeHTTP(w, r)
	}
}
