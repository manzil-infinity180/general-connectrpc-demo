package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OptionRepo struct {
	db *pgxpool.Pool
}

func NewOptionRepo(db *pgxpool.Pool) *OptionRepo {
	return &OptionRepo{db: db}
}

func (r *OptionRepo) CreateOption(ctx context.Context, pollID, text string) error {
	query := `INSERT INTO options (poll_id, text) VALUES ($1,$2)`
	_, err := r.db.Exec(ctx, query, pollID, text)
	return err
}

func (r *OptionRepo) IncrementVoteCount(ctx context.Context, optionID string) error {
	query := `UPDATE options SET vote_count = vote_count + 1 WHERE id=$1`
	_, err := r.db.Exec(ctx, query, optionID)
	return err
}
