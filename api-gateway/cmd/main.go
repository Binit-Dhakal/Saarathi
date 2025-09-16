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
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// might need to add checkOrigin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	userServiceURL, _ := url.Parse("http://users-service:8080")
	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	tripServiceURL, _ := url.Parse("http://trips-service:8082")
	tripServiceProxy := httputil.NewSingleHostReverseProxy(tripServiceURL)

	driverStateURL, _ := url.Parse("http://driver-state-service:8084")
	driverStateProxy := httputil.NewSingleHostReverseProxy(driverStateURL)
	driverStateProxy.Director = func(req *http.Request) {
		req.URL.Scheme = driverStateURL.Scheme
		req.URL.Host = driverStateURL.Host
		req.URL.Path = "/ws"
	}
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

	mux.Handle("/ws/driver", authMiddleware(proxyHandler(driverStateProxy)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      rest.CorsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprint(os.Stderr, "server failed", err)
		os.Exit(1)
	}
}

func proxyHandler(p *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
