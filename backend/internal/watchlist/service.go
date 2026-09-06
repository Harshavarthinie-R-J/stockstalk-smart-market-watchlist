package watchlist

import (
	"context"
	"errors"
	"strings"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, userID, name string) (*Watchlist, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("watchlist name is required")
	}
	return s.repo.Create(ctx, userID, name)
}

func (s *Service) Get(ctx context.Context, userID, id string) (*Watchlist, error) {
	return s.repo.GetByID(ctx, userID, id)
}
func (s *Service) List(ctx context.Context, userID string) ([]Watchlist, error) {
	return s.repo.List(ctx, userID)
}
func (s *Service) AddStock(ctx context.Context, userID, watchlistID, instrumentID string) error {
	return s.repo.AddStock(ctx, userID, watchlistID, instrumentID)
}
func (s *Service) RemoveStock(ctx context.Context, userID, watchlistID, instrumentID string) error {
	return s.repo.RemoveStock(ctx, userID, watchlistID, instrumentID)
}

// ListAll returns all watchlists.
// Used by the background market processor.
func (s *Service) ListAll(ctx context.Context) ([]Watchlist, error) {
	return s.repo.ListAll(ctx)
}
