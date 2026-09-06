package instrument

import "context"

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }
func (s *Service) GetByID(ctx context.Context, id string) (*Instrument, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) GetBySymbol(ctx context.Context, symbol string) (*Instrument, error) {
	return s.repo.GetBySymbol(ctx, symbol)
}
func (s *Service) Search(ctx context.Context, query string) ([]Instrument, error) {
	return s.repo.Search(ctx, query)
}
