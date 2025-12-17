package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	votingv1 "rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1"
	"rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1/votingv1connect"
	"time"

	"connectrpc.com/connect"
)

func main() {
	ctx := context.Background()

	// HTTP client (ConnectRPC uses HTTP/2 automatically)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// ConnectRPC client
	client := votingv1connect.NewVotingServiceClient(
		httpClient,
		"http://localhost:8080",
	)

	// ---- 1️⃣ Create Poll ----
	fmt.Println("Creating poll...")

	createReq := connect.NewRequest(&votingv1.PollCreationRequest{
		Question:    "Best backend language?",
		OptionTexts: []string{"Go", "Rust", "Java"},
	})

	// Fake Google ID token
	createReq.Header().Set("Authorization", "Bearer FAKE_GOOGLE_ID_TOKEN")

	createResp, err := client.CreatePoll(ctx, createReq)
	if err != nil {
		log.Fatalf("CreatePoll failed: %v", err)
	}

	pollID := createResp.Msg.PollId
	fmt.Println("Poll created:", pollID)

	// ---- 2️⃣ Submit Vote ----
	fmt.Println("Submitting vote...")

	voteReq := connect.NewRequest(&votingv1.VoteSubmission{
		PollId:      pollID,
		OptionId:    "OPTION_ID_FROM_DB", // manually copy from DB for now
		VoterIdHash: "fake-browser-fingerprint",
	})

	voteResp, err := client.SubmitVote(ctx, voteReq)
	if err != nil {
		log.Fatalf("SubmitVote failed: %v", err)
	}

	fmt.Println("Vote result:", voteResp.Msg.Success)

	// ---- 3️⃣ Stream Results (optional) ----
	fmt.Println("Streaming results...")

	streamReq := connect.NewRequest(&votingv1.PollRequest{
		PollId: pollID,
	})

	stream, err := client.StreamResults(ctx, streamReq)
	if err != nil {
		log.Fatalf("StreamResults failed: %v", err)
	}

	go func() {
		for stream.Receive() {
			update := stream.Msg()
			fmt.Println("🔴 Stream update:", update.PollId)
		}
		if err := stream.Err(); err != nil {
			log.Println("stream error:", err)
		}
	}()

	// Keep client alive
	time.Sleep(30 * time.Second)
}
