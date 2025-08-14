package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/api-gateway/internal/handlers/rest"
)

func main() {
	userServiceURL, _ := url.Parse("http://users-service:8080")
	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/users/riders", rest.CorsMiddleware(proxyHandler(userServiceProxy)))
	mux.Handle("/api/v1/tokens", rest.CorsMiddleware(proxyHandler(userServiceProxy)))

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("starting server on :8081")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprint(os.Stderr, "server failed", err)
		os.Exit(1)
	}
}

func proxyHandler(p *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Forwarding request to %s", r.URL.Path)
		p.ServeHTTP(w, r)
	}
}

func getPublicKey() (*rsa.PublicKey, error) {
	keyString := os.Getenv("JWT_PUBLIC_KEY")
	if keyString == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY not set")
	}

	block, _ := pem.Decode([]byte(keyString))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}

	return rsaKey, nil
}
