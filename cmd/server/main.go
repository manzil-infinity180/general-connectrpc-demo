package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"rahulxf.com/general-connectrpc-demo/internal/auth"
	"rahulxf.com/general-connectrpc-demo/internal/db"
	votingv1 "rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1"
	"rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1/votingv1connect"
	"rahulxf.com/general-connectrpc-demo/internal/repository"
	"rahulxf.com/general-connectrpc-demo/internal/service"
	"rahulxf.com/general-connectrpc-demo/internal/stream"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	log.Println("Using DATABASE_URL:", dsn)
	dbPool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// streaming hub
	hub := stream.NewHub[*votingv1.PollResultUpdate]()

	// repositories
	pollRepo := repository.NewPollRepo(dbPool)
	optionRepo := repository.NewOptionRepo(dbPool)
	voteRepo := repository.NewVoteRepo(dbPool)

	// service
	votingService := service.NewVotingService(
		pollRepo,
		optionRepo,
		voteRepo,
		hub,
	)

	// mux
	mux := http.NewServeMux()
	mux.Handle(
		votingv1connect.NewVotingServiceHandler(votingService),
	)

	userRepo := repository.NewUserRepo(dbPool)
	handler := auth.DevMiddleware(userRepo)(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("🚀 Voting service running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
