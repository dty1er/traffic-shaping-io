package tokenbucket

import (
	"context"
	"fmt"
	"time"
)

type Store interface {
	GetState(ctx context.Context) (*State, error)
	SetState(ctx context.Context, state *State) error
}

type Bucket struct {
	// rps is a rate that tokens are filled in the bucket.
	// A token is refilled every 1/rate second.
	// Bigger rate, faster tokens refillment.
	rps int64

	// capacity
	capacity int64

	store Store
}

type State struct {
	Last            int64 // unix timestamp (nanosec)
	AvailableTokens int64 // number of available tokens
}

func (s *State) IsZero() bool {
	return s.Last == 0
}

func New(rps, capacity int64, store Store) *Bucket {
	return &Bucket{rps: rps, capacity: capacity, store: store}
}

func (b *Bucket) Take(ctx context.Context, n int) (time.Duration, error) {
	now := time.Now().Unix() / int64(time.Millisecond)

	prev, err := b.store.GetState(ctx)
	if err != nil {
		return 0, err
	}

	newState := &State{Last: now}

	if prev.IsZero() {
		// comes here only on first call
		prev.Last = newState.Last
		prev.AvailableTokens = b.capacity
	}

	sub := now - prev.Last                // milliseconds
	refill := sub / (time.Second / b.rps) // how many token should be refilled
	prev.AvailableTokens += refill
	if prev.AvailableTokens > b.capacity {
		prev.AvailableTokens = b.capacity
	}

	prev.AvailableTokens -= n
	if prev.AvailableTokens < 0 {
		// todo: sleep
		prev.AvailableTokens = 0
	}

	newState.AvailableTokens = prev.AvailableTokens
	// need to lock
	b.store.SetState(ctx, newState)
}

func (b *Bucket) TakeOne(ctx context.Context) (time.Duration, error) {
	return b.Take(ctx, 1)
}

func (b *Bucket) Debug() {
	fmt.Printf("rate: %v, capacity: %v, num of tokens: %v\n", b.rate, b.capacity, len(b.tokens))
}
