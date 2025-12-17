package main

import (
	"context"
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

	db, _ := db.New(ctx, os.Getenv("DATABASE_URL"))

	hub := stream.NewHub[*votingv1.PollResultUpdate]()

	svc := service.NewVotingService(
		repository.NewPollRepo(db),
		repository.NewOptionRepo(db),
		repository.NewVoteRepo(db),
		hub,
	)

	mux := http.NewServeMux()
	mux.Handle(votingv1connect.NewVotingServiceHandler(svc))

	http.ListenAndServe(":8080", auth.Middleware(mux))
}
