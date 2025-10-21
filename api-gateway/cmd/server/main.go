package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	authmw "api-gateway/internal/middleware"
)

func main() {
	// Get service addresses from environment variables
	// In production, these would come from service discovery (Consul, K8s DNS, etc.)
	authServiceAddr := getEnv("AUTH_SERVICE_URL", "localhost:50051")
	userServiceAddr := getEnv("USER_SERVICE_URL", "localhost:50052")
	port := getEnv("PORT", "8080")

	// Connect to all backend gRPC services
	log.Println("Connecting to backend services...")
	grpcClients, err := clients.NewGRPCClients(authServiceAddr, userServiceAddr)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC services: %v", err)
	}
	defer grpcClients.Close() // Ensure connections are closed on shutdown

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(grpcClients.AuthClient)
	userHandler := handlers.NewUserHandler(grpcClients.UserClient)

	// Setup router with Chi
	r := chi.NewRouter()

	// Global middleware applies to ALL routes
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check endpoint (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes under /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public - no authentication required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// User routes (protected - require authentication)
		r.Route("/users", func(r chi.Router) {
			// Apply auth middleware to all user routes
			r.Use(authmw.AuthMiddleware(grpcClients.AuthClient))

			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)

			// Address sub-routes
			r.Post("/{id}/addresses", userHandler.AddAddress)
			r.Get("/{id}/addresses", userHandler.GetAddresses)
		})

		// TODO: Add product, order, payment routes as you build those services
	})

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine so it doesn't block shutdown handling
	go func() {
		log.Printf("🚀 API Gateway listening on port %s", port)
		log.Printf("   Auth Service: %s", authServiceAddr)
		log.Printf("   User Service: %s", userServiceAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown handling
	// This is important for production - allows in-flight requests to complete
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until we receive a signal

	log.Println("Shutting down API Gateway...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway stopped gracefully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
