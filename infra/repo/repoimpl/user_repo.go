package repoimpl

import (
	"context"

	"gorm.io/gorm"

	"funny-project-be/domain/entity"
)

// UserRepo implements methods of user's repository.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates and returns a new instances of UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Get finds and returns a user by id.
func (r *UserRepo) Get(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User

	query := r.db

	if err := query.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOneByEmail finds and returns a user by email.
func (r *UserRepo) GetOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	query := r.db

	if err := query.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Add adds new users to repo.
func (r *UserRepo) Add(ctx context.Context, users ...*entity.User) error {
	for _, user := range users {
		if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}
	}

	return nil
}

// Remove removes users from repo.
func (r *UserRepo) Remove(ctx context.Context, users ...*entity.User) error {
	for _, user := range users {
		if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update updates users in repo.
func (r *UserRepo) Update(ctx context.Context, users ...*entity.User) error {
	for _, user := range users {
		if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
			return err
		}
	}

	return nil
}
