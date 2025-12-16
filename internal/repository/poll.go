package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PollRepo struct {
	db *pgxpool.Pool
}

func NewPollRepo(db *pgxpool.Pool) *PollRepo {
	return &PollRepo{db: db}
}

func (r *PollRepo) CreatePoll(ctx context.Context, question, creatorID string) (string, error) {
	var pollID string
	query := `INSERT INTO polls (question, creator_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, question, creatorID).Scan(&pollID)
	return pollID, err
}

func (r *PollRepo) ClosePoll(ctx context.Context, pollID, userID string) error {
	query := `UPDATE polls SET status='CLOSED', closed_at=now() WHERE id=$1 AND creator_id=$2`
	_, err := r.db.Exec(ctx, query, pollID, userID)
	return err
}
