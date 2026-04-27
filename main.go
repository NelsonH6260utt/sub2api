package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appName        = "sub2api"
	appVersion     = "1.0.0"
)

func main() {
	var (
		host    string
		port    int
		version bool
	)

	flag.StringVar(&host, "host", getEnvOrDefault("HOST", defaultHost), "Host address to listen on")
	flag.IntVar(&port, "port", getEnvIntOrDefault("PORT", defaultPort), "Port to listen on")
	flag.BoolVar(&version, "version", false, "Print version information and exit")
	flag.Parse()

	if version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthCheck)

	// Subscription conversion endpoints
	mux.HandleFunc("/sub", handler.ConvertSubscription)
	mux.HandleFunc("/api/v1/convert", handler.ConvertSubscription)

	log.Printf("Starting %s v%s on %s", appName, appVersion, addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value.
func getEnvOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// getEnvIntOrDefault returns the integer value of an environment variable or a default value.
func getEnvIntOrDefault(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
