package auth

import (
	"context"
	"net/http"

	"rahulxf.com/general-connectrpc-demo/internal/repository"
)

func DevMiddleware(userRepo *repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Only authenticate if Authorization header exists
			if r.Header.Get("Authorization") != "" {

				userID, err := EnsureDevUser(r.Context(), userRepo)
				if err != nil {
					http.Error(w, "failed to create dev user", http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
