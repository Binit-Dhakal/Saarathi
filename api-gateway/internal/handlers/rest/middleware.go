package rest

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Binit-Dhakal/Saarathi/pkg/cookies"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Binit-Dhakal/Saarathi/pkg/claims"
)

type MiddlewareFunc func(http.Handler) http.Handler

func CorsMiddleware(next http.Handler) http.Handler {
	trustedOriginsStr := os.Getenv("TRUSTED_ORIGINS")
	trustedOrigins := strings.Split(trustedOriginsStr, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

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
			fmt.Println(r.Cookies())
			authToken, err := cookies.Read(r, "accessToken")
			fmt.Println(authToken, err)
			if err != nil {
				http.Error(w, "Cookie Read error", http.StatusBadRequest)
				return
			}

			token, err := jwt.ParseWithClaims(authToken, &claims.CustomClaims{}, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*claims.CustomClaims)
			if !ok {
				http.Error(w, "Invalid Claims", http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			next.ServeHTTP(w, r)
		})
	}
}
