package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
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

	tripServiceURL, _ := url.Parse("http://trips-service:8070")
	tripServiceProxy := httputil.NewSingleHostReverseProxy(tripServiceURL)

	riderServiceURL, _ := url.Parse("http://rider-service:8010")
	riderServiceProxy := httputil.NewSingleHostReverseProxy(riderServiceURL)

	riderServiceProxy.FlushInterval = -1
	riderServiceProxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false,
		DisableCompression:    true,
	}

	driverStateURL, _ := url.Parse("http://driver-state-service:8050")
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

	mux.Handle("/api/v1/trip/", authMiddleware(proxyHandler(riderServiceProxy)))

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
		fmt.Println("Serving: ", r.URL)
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
