package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	pb "go-project/proto/auth"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://authuser:authpass@localhost:5432/authdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "23094uoifhpowisrfhp02893ryufpowafh"
	}

	userRepo := repository.NewPostgresUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Auth service listening on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
