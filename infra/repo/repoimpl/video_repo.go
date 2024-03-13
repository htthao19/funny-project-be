package repoimpl

import (
	"context"
	"errors"

	"funny-project-be/domain/entity"

	"gorm.io/gorm"
)

// VideoRepo implements methods of video's repository.
type VideoRepo struct {
	db *gorm.DB
}

// NewVideoRepo creates and returns a new instances of VideoRepo.
func NewVideoRepo(db *gorm.DB) *VideoRepo {
	return &VideoRepo{db: db}
}

// Get finds and returns a video by id.
func (r *VideoRepo) Get(ctx context.Context, id uint) (*entity.Video, error) {
	var video entity.Video

	query := r.db

	if err := query.First(&video, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

// GetRangeByQuery finds and returns a range of videos.
func (r *VideoRepo) GetRangeByQuery(ctx context.Context, sort string, limit int, page int) ([]*entity.Video, error) {
	var videos []*entity.Video

	q := r.db

	if sort != "" {
		q = q.Order(sort)
	} else {
		q = q.Order("id desc")
	}

	if err := q.Limit(limit).Offset(limit * (page - 1)).Find(&videos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*entity.Video{}, nil
		}
		return nil, err
	}

	return videos, nil
}

// Count counts and returns the number of videos.
func (r *VideoRepo) Count(ctx context.Context) (int64, error) {
	var videos []*entity.Video
	var count int64
	q := r.db
	if err := q.Find(&videos).Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return -1, err
	}

	return count, nil
}

// Add adds new videos to repo.
func (r *VideoRepo) Add(ctx context.Context, videos ...*entity.Video) error {
	for _, video := range videos {
		if err := r.db.WithContext(ctx).Create(video).Error; err != nil {
			return err
		}
	}

	return nil
}

// Remove removes videos from repo.
func (r *VideoRepo) Remove(ctx context.Context, videos ...*entity.Video) error {
	for _, video := range videos {
		if err := r.db.WithContext(ctx).Delete(video).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update updates videos in repo.
func (r *VideoRepo) Update(ctx context.Context, videos ...*entity.Video) error {
	for _, video := range videos {
		if err := r.db.WithContext(ctx).Save(video).Error; err != nil {
			return err
		}
	}

	return nil
}
