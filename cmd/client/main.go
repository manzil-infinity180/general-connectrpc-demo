package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	votingv1 "rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1"
	"rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1/votingv1connect"
)

func main() {
	ctx := context.Background()

	fmt.Println("Voting System Demo Client")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	client := votingv1connect.NewVotingServiceClient(
		httpClient,
		"http://localhost:8080",
	)

	fmt.Println("STEP 1: Creating Poll")

	createReq := connect.NewRequest(&votingv1.PollCreationRequest{
		Question:    "What's the best backend language in 2024?",
		OptionTexts: []string{"Go", "Rust", "Java", "Python"},
	})
	createReq.Header().Set("Authorization", "Bearer fake-dev-token")

	createResp, err := client.CreatePoll(ctx, createReq)
	if err != nil {
		log.Fatalf("CreatePoll failed: %v", err)
	}
	pollID := createResp.Msg.PollId
	fmt.Printf("Poll created successfully!\n")
	fmt.Printf("Poll ID: %s\n\n", pollID)

	fmt.Println("STEP 2: Fetching Options from Database")

	type Option struct {
		ID   string
		Text string
	}
	query := `SELECT id, text FROM options WHERE poll_id = $1 ORDER BY text`
	rows, err := dbPool.Query(ctx, query, pollID)
	if err != nil {
		log.Fatalf("Failed to query options: %v", err)
	}

	var options []Option
	for rows.Next() {
		var opt Option
		if err := rows.Scan(&opt.ID, &opt.Text); err != nil {
			log.Fatalf("Failed to scan option: %v", err)
		}
		options = append(options, opt)
	}
	rows.Close()

	if len(options) == 0 {
		log.Fatal("No options found for poll")
	}

	fmt.Printf("Found %d options:\n", len(options))
	for i, opt := range options {
		fmt.Printf("   %d. %s (ID: %s)\n", i+1, opt.Text, opt.ID)
	}

	fmt.Println("STEP 3: Starting Real-time Stream")

	streamReq := connect.NewRequest(&votingv1.PollRequest{
		PollId: pollID,
	})

	stream, err := client.StreamResults(ctx, streamReq)
	if err != nil {
		log.Fatalf("StreamResults failed: %v", err)
	}
	streamActive := true
	go func() {
		updateCount := 0
		for stream.Receive() {
			update := stream.Msg()
			updateCount++

			fmt.Printf("\nLIVE UPDATE #%d\n", updateCount)
			fmt.Printf("Poll: %s\n", update.PollId)

			if len(update.UpdatedOptions) > 0 {
				fmt.Println("   Current Results:")
				for _, opt := range update.UpdatedOptions {
					fmt.Printf("     • %s: %d vote(s)\n", opt.OptionText, opt.VoteCount)
				}
				fmt.Printf("   Total Votes: %d\n", update.TotalVotesCast)
			}
			fmt.Println()
		}

		if err := stream.Err(); err != nil {
			log.Printf("Stream ended: %v\n", err)
		}
		streamActive = false
	}()

	fmt.Println("Stream connected! Waiting for vote updates...\n")
	time.Sleep(1 * time.Second)

	fmt.Println("STEP 4: Submitting Votes")

	// Vote 1: Vote for "Go"
	goOption := options[0] // First option (Go)
	fmt.Printf("Vote #1: Voting for '%s'...\n", goOption.Text)

	voteReq1 := connect.NewRequest(&votingv1.VoteSubmission{
		PollId:      pollID,
		OptionId:    goOption.ID,
		VoterIdHash: "voter-alice-12345",
	})

	voteResp1, err := client.SubmitVote(ctx, voteReq1)
	if err != nil {
		fmt.Printf("Vote failed: %v\n", err)
	} else {
		fmt.Printf("Vote recorded! (Success: %v)\n", voteResp1.Msg.Success)
	}

	time.Sleep(2 * time.Second)

	// Vote 2: Vote for "Rust"
	if len(options) > 1 {
		rustOption := options[1]
		fmt.Printf("Vote #2: Voting for '%s'...\n", rustOption.Text)

		voteReq2 := connect.NewRequest(&votingv1.VoteSubmission{
			PollId:      pollID,
			OptionId:    rustOption.ID,
			VoterIdHash: "voter-bob-67890",
		})

		voteResp2, err := client.SubmitVote(ctx, voteReq2)
		if err != nil {
			fmt.Printf("Vote failed: %v\n", err)
		} else {
			fmt.Printf("Vote recorded! (Success: %v)\n", voteResp2.Msg.Success)
		}

		time.Sleep(2 * time.Second)
	}

	// Vote 3: Try duplicate vote (should fail)
	fmt.Printf("Vote #3: Trying to vote again with same voter (should fail)...\n")

	duplicateReq := connect.NewRequest(&votingv1.VoteSubmission{
		PollId:      pollID,
		OptionId:    goOption.ID,
		VoterIdHash: "voter-alice-12345", // Same voter as before
	})

	_, err = client.SubmitVote(ctx, duplicateReq)
	if err != nil {
		fmt.Printf("Duplicate vote prevented! Error: %v\n", err)
	} else {
		fmt.Println("Warning: Duplicate vote was allowed (should be prevented)")
	}

	time.Sleep(2 * time.Second)

	// Vote 4: Vote for "Java"
	if len(options) > 2 {
		javaOption := options[2]
		fmt.Printf("Vote #4: Voting for '%s'...\n", javaOption.Text)

		voteReq4 := connect.NewRequest(&votingv1.VoteSubmission{
			PollId:      pollID,
			OptionId:    javaOption.ID,
			VoterIdHash: "voter-charlie-11111",
		})

		voteResp4, err := client.SubmitVote(ctx, voteReq4)
		if err != nil {
			fmt.Printf("Vote failed: %v\n", err)
		} else {
			fmt.Printf("Vote recorded! (Success: %v)\n", voteResp4.Msg.Success)
		}
	}

	fmt.Println("\nSTEP 5: Final Results")
	fmt.Println("────────────────────────")

	time.Sleep(2 * time.Second)

	query = `
        SELECT o.text, o.vote_count 
        FROM options o 
        WHERE o.poll_id = $1 
        ORDER BY o.vote_count DESC, o.text
    `
	rows, err = dbPool.Query(ctx, query, pollID)
	if err != nil {
		log.Printf("Failed to query final results: %v", err)
	} else {
		fmt.Println("Final Vote Counts:")
		totalVotes := 0
		for rows.Next() {
			var text string
			var count int
			if err := rows.Scan(&text, &count); err != nil {
				continue
			}
			fmt.Printf("   • %s: %d vote(s)\n", text, count)
			totalVotes += count
		}
		fmt.Printf("\n   Total: %d vote(s) cast\n", totalVotes)
		rows.Close()
	}

	fmt.Println("\nSTEP 6: Closing Poll")

	closeReq := connect.NewRequest(&votingv1.PollRequest{
		PollId: pollID,
	})
	closeReq.Header().Set("Authorization", "Bearer fake-dev-token")

	_, err = client.ClosePoll(ctx, closeReq)
	if err != nil {
		fmt.Printf("Failed to close poll: %v\n", err)
	} else {
		fmt.Println("Poll closed successfully!")
	}

	time.Sleep(1 * time.Second)

	closedVoteReq := connect.NewRequest(&votingv1.VoteSubmission{
		PollId:      pollID,
		OptionId:    goOption.ID,
		VoterIdHash: "voter-late-99999",
	})

	_, err = client.SubmitVote(ctx, closedVoteReq)
	if err != nil {
		fmt.Printf("Vote prevented on closed poll! Error: %v\n", err)
	} else {
		fmt.Println("Warning: Vote on closed poll was allowed (should be prevented)")
	}

	fmt.Printf("\nPoll ID: %s\n", pollID)
	fmt.Println("\nYou can view this poll in the database:")
	fmt.Printf("psql $DATABASE_URL -c \"SELECT * FROM polls WHERE id='%s';\"\n", pollID)

	if streamActive {
		fmt.Println("\nStream still active. Press Ctrl+C to exit.")
		time.Sleep(10 * time.Second)
	}
}
