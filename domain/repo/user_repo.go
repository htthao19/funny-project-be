package repo

import (
	"context"

	"funny-project-be/domain/entity"
)

// UserRepo exposes methods of user's repository.
type UserRepo interface {
	// Get finds and returns a user by id.
	Get(ctx context.Context, id uint) (*entity.User, error)

	// GetOneByEmail finds and returns a user by email.
	GetOneByEmail(ctx context.Context, email string) (*entity.User, error)

	// Add adds new users to repo.
	Add(ctx context.Context, users ...*entity.User) error

	// Remove removes users from repo.
	Remove(ctx context.Context, users ...*entity.User) error

	// Update updates users in repo.
	Update(ctx context.Context, users ...*entity.User) error
}
