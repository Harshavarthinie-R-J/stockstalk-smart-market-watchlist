package graphql

import (
	"context"
	"sync"
)

type eventBroker struct {
	mu   sync.RWMutex
	subs map[string]map[chan interface{}]struct{}
}

func newEventBroker() *eventBroker {
	return &eventBroker{
		subs: make(map[string]map[chan interface{}]struct{}),
	}
}

func subscriptionKey(userID, watchlistID string) string {
	return userID + ":" + watchlistID
}

func (b *eventBroker) Subscribe(
	ctx context.Context,
	userID string,
	watchlistID string,
) <-chan interface{} {
	ch := make(chan interface{}, 10)
	key := subscriptionKey(userID, watchlistID)

	b.mu.Lock()

	if b.subs[key] == nil {
		b.subs[key] = make(map[chan interface{}]struct{})
	}

	b.subs[key][ch] = struct{}{}

	b.mu.Unlock()

	go func() {
		<-ctx.Done()

		b.mu.Lock()

		if subscribers, ok := b.subs[key]; ok {
			delete(subscribers, ch)

			if len(subscribers) == 0 {
				delete(b.subs, key)
			}
		}

		close(ch)

		b.mu.Unlock()
	}()

	return ch
}

func (b *eventBroker) Publish(
	userID string,
	watchlistID string,
	event interface{},
) {
	key := subscriptionKey(userID, watchlistID)

	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subs[key] {
		select {
		case ch <- event:
		default:
			// Do not allow a slow subscriber to block
			// the market-change processing pipeline.
		}
	}
}