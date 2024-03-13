package repo

import (
	"context"

	"funny-project-be/domain/entity"
)

// VideoRepo exposes methods of video's repository.
type VideoRepo interface {
	// Get finds and returns a video by id.
	Get(ctx context.Context, id uint) (*entity.Video, error)

	// GetRangeByQuery finds and returns a range of videos.
	GetRangeByQuery(ctx context.Context, sort string, limit int, page int) ([]*entity.Video, error)

	// Count counts and returns the number of videos.
	Count(ctx context.Context) (int64, error)

	// Add adds new videos to repo.
	Add(ctx context.Context, videos ...*entity.Video) error

	// Remove removes videos from repo.
	Remove(ctx context.Context, videos ...*entity.Video) error

	// Update updates videos in repo.
	Update(ctx context.Context, videos ...*entity.Video) error
}
