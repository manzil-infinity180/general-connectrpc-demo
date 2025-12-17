package service

import (
	"connectrpc.com/connect"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"google.golang.org/protobuf/types/known/emptypb"
	"rahulxf.com/general-connectrpc-demo/internal/auth"
	votingv1 "rahulxf.com/general-connectrpc-demo/internal/gen/go/voting/v1"
	"rahulxf.com/general-connectrpc-demo/internal/repository"
	"rahulxf.com/general-connectrpc-demo/internal/stream"
)

type VotingService struct {
	polls   *repository.PollRepo
	options *repository.OptionRepo
	votes   *repository.VoteRepo
	hub     *stream.Hub[*votingv1.PollResultUpdate]
}

func NewVotingService(
	polls *repository.PollRepo,
	options *repository.OptionRepo,
	votes *repository.VoteRepo,
	hub *stream.Hub[*votingv1.PollResultUpdate],
) *VotingService {
	return &VotingService{
		polls:   polls,
		options: options,
		votes:   votes,
		hub:     hub,
	}
}

func (s *VotingService) CreatePoll(ctx context.Context,
	req *connect.Request[votingv1.PollCreationRequest]) (*connect.Response[votingv1.PollCreationResponse], error) {

	userID := ctx.Value(auth.UserIDKey).(string)
	pollID, _ := s.polls.CreatePoll(ctx, req.Msg.Question, userID)
	for _, opt := range req.Msg.OptionTexts {
		s.options.CreateOption(ctx, pollID, opt)
	}
	return connect.NewResponse(&votingv1.PollCreationResponse{
		PollId: pollID,
	}), nil
}

func hashVoter(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}

func (s *VotingService) SubmitVote(ctx context.Context,
	req *connect.Request[votingv1.VoteSubmission]) (*connect.Response[votingv1.VoteReceipt], error) {

	voterHash := hashVoter(req.Msg.VoterIdHash)
	s.votes.InsertVote(ctx, req.Msg.PollId, req.Msg.OptionId, voterHash)
	s.options.IncrementVoteCount(ctx, req.Msg.OptionId)

	s.hub.Publish(req.Msg.PollId, &votingv1.PollResultUpdate{
		PollId: req.Msg.PollId,
	})

	return connect.NewResponse(&votingv1.VoteReceipt{
		Success: true,
	}), nil
}

func (s *VotingService) StreamResults(
	ctx context.Context,
	req *connect.Request[votingv1.PollRequest],
	stream *connect.ServerStream[votingv1.PollResultUpdate],
) error {

	ch := s.hub.Subscribe(req.Msg.PollId)

	for {
		select {
		case msg := <-ch:
			stream.Send(msg)
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *VotingService) ClosePoll(
	ctx context.Context,
	req *connect.Request[votingv1.PollRequest],
) (*connect.Response[emptypb.Empty], error) {

	userID, _ := ctx.Value(auth.UserIDKey).(string)

	if err := s.polls.ClosePoll(ctx, req.Msg.PollId, userID); err != nil {
		return nil, err
	}

	// notify stream listeners poll is closed
	s.hub.Publish(req.Msg.PollId, &votingv1.PollResultUpdate{
		PollId: req.Msg.PollId,
	})

	return connect.NewResponse(&emptypb.Empty{}), nil
}
