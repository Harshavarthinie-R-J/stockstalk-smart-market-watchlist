package instrument

import (
	"context"
	"errors"
	"strings"
	"sync"
)

type Repository struct {
	mu          sync.RWMutex
	instruments map[string]Instrument
}

func NewRepository(instruments []Instrument) *Repository {
	repo := &Repository{
		instruments: make(map[string]Instrument),
	}

	_ = repo.UpsertMany(context.Background(), instruments)

	return repo
}

// Store instruments.
func (r *Repository) UpsertMany(ctx context.Context, items []Instrument) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		r.instruments[item.ID] = item
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Instrument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.instruments[id]
	if !ok {
		return nil, errors.New("instrument not found")
	}

	result := item
	return &result, nil
}

func (r *Repository) GetBySymbol(ctx context.Context, symbol string) (*Instrument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.instruments {
		if strings.EqualFold(item.Symbol, symbol) {
			result := item
			return &result, nil
		}
	}

	return nil, errors.New("instrument not found")
}

func (r *Repository) Search(ctx context.Context, query string) ([]Instrument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))

	result := make([]Instrument, 0)

	for _, item := range r.instruments {
		if !item.Active {
			continue
		}

		if query == "" ||
			strings.Contains(strings.ToLower(item.Symbol), query) ||
			strings.Contains(strings.ToLower(item.Name), query) {
			result = append(result, item)
		}
	}

	return result, nil
}
