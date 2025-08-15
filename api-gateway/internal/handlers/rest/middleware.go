package rest

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Binit-Dhakal/Saarathi/pkg/cookies"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareFunc func(http.Handler) http.Handler

func CorsMiddleware(next http.Handler) http.Handler {
	trustedOriginsStr := os.Getenv("TRUSTED_ORIGINS")
	trustedOrigins := strings.Split(trustedOriginsStr, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if origin != "" {
			for i := range trustedOrigins {
				if origin == trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// this is pre-flight CORS request
						w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func NewAuthMiddleware(publicKey *rsa.PublicKey) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken, err := cookies.Read(r, "authToken")
			if err != nil {
				http.Error(w, "Cookie Read error", http.StatusBadRequest)
				return
			}

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})

			if err != nil {
				log.Printf("Token validation failed: %v", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				log.Println("Received an invalid token")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		})
	}
}
