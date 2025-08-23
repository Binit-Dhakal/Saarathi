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
	"strings"
	"sync"
	"time"

	"github.com/Binit-Dhakal/Saarathi/api-gateway/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/pkg/env"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// might need to add checkOrigin
}

func main() {
	userServiceURL, _ := url.Parse("http://users-service:8080")
	userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	tripServiceURL, _ := url.Parse("http://trips-service:8082")
	tripServiceProxy := httputil.NewSingleHostReverseProxy(tripServiceURL)

	// driverWsURL, _ := url.Parse("ws://driver-state-service:8084/ws")
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
		p.ServeHTTP(w, r)
	}
}

func websocketProxyHandler(target *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade client connection: %v", err)
			return
		}

		headers := make(http.Header)
		for key, values := range r.Header {
			// Skip WebSocket-specific headers
			if strings.HasPrefix(key, "Sec") || key == "Upgrade" || key == "Connection" {
				continue
			}
			for _, value := range values {
				headers.Add(key, value)
			}
		}

		backendConn, _, err := websocket.DefaultDialer.Dial(target.String(), headers)
		if err != nil {
			log.Printf("Failed to dial to backend service: %v", err)
			clientConn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseServiceRestart,
					"Backend service unavailable",
				),
			)
			return
		}

		var once sync.Once
		closeConns := func() {
			backendConn.Close()
			clientConn.Close()
		}

		// proxy data from client to backend
		go func() {
			defer once.Do(closeConns)
			for {
				messageType, p, err := clientConn.ReadMessage()
				if err != nil {
					log.Printf("client read error: %v", err)
					return
				}

				err = backendConn.WriteMessage(messageType, p)
				if err != nil {
					log.Printf("Backend write error: %v", err)
					return
				}
			}
		}()

		// proxy data from backend to client
		go func() {
			once.Do(closeConns)
			for {
				messageType, p, err := backendConn.ReadMessage()
				if err != nil {
					log.Printf("backend read error: %v", err)
					return
				}

				err = clientConn.WriteMessage(messageType, p)
				if err != nil {
					log.Printf("Client write error: %v", err)
					return
				}

			}
		}()

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
