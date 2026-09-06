package checkpoint

import "context"

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, userID, watchlistID string) (*Checkpoint, error) {
	return s.repo.Create(ctx, userID, watchlistID)
}
func (s *Service) Latest(ctx context.Context, userID, watchlistID string) (*Checkpoint, error) {
	return s.repo.Latest(ctx, userID, watchlistID)
}
