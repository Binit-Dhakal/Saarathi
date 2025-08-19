package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/api-gateway/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/pkg/env"
)

func main() {
	userServiceURL, _ := url.Parse("http://users-service:8080")
	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	tripServiceURL, _ := url.Parse("http://trips-service:8082")
	tripServiceProxy := httputil.NewSingleHostReverseProxy(tripServiceURL)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/users/", proxyHandler(userServiceProxy))
	mux.Handle("/api/v1/tokens/", proxyHandler(userServiceProxy))

	// Authenticated request
	publicKey, err := getPublicKey()
	if err != nil {
		panic(err)
	}

	authMiddleware := rest.NewAuthMiddleware(publicKey)
	mux.Handle("/api/v1/fare/", authMiddleware(proxyHandler(tripServiceProxy)))

	server := &http.Server{
		Addr:         ":8081",
		Handler:      rest.CorsMiddleware(mux),
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
		fmt.Printf("Serving: %v", r.URL)
		p.ServeHTTP(w, r)
	}
}

func getPublicKey() (*rsa.PublicKey, error) {
	keyString, err := env.GetEnv("JWT_PUBLIC_KEY")
	if err != nil {
		return nil, err
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
