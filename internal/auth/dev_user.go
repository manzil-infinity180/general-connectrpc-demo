package auth

import (
	"context"

	"rahulxf.com/general-connectrpc-demo/internal/model"
	"rahulxf.com/general-connectrpc-demo/internal/repository"
)

func EnsureDevUser(ctx context.Context, userRepo *repository.UserRepo) (string, error) {
	user := &model.User{
		Email:     "dev@local.test",
		Name:      "Local Dev",
		GoogleSub: "local-dev-google-sub",
	}

	dbUser, err := userRepo.GetOrCreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return dbUser.ID, nil
}
