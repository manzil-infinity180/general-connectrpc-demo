package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepo struct {
	db *pgxpool.Pool
}

func NewVoteRepo(db *pgxpool.Pool) *VoteRepo {
	return &VoteRepo{db: db}
}

func (r *VoteRepo) InsertVote(ctx context.Context, pollID, optionID, voterHash string) error {
	query := `INSERT INTO votes (poll_id, option_id, voter_hash) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(ctx, query, pollID, optionID, voterHash)
	return err
}
