package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"rahulxf.com/general-connectrpc-demo/internal/model"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOrCreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	query := `INSERT INTO users (email, name, google_sub)
			VALUES ($1, $2, $3)
			ON CONFLICT (google_sub)
			DO UPDATE SET email=EXCLUDED.email
			RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, u.Email, u.Name, u.GoogleSub).Scan(&u.ID, &u.CreatedAt)
	return u, err
}
