package tokenbucket

import (
	"fmt"
	"time"
)

type Bucket struct {
	// rate is a rate that tokens are filled in the bucket.
	// A token is refilled every 1/rate second.
	// Bigger rate, faster tokens refillment.
	rate int64

	// capacity is a capacity of this bucket.
	// It is the same as upper boundary of number of tokens in the bucket
	capacity int64

	// tokens represents tokens in this bucket
	tokens chan struct{}

	done chan struct{}
}

func New(rate, capacity int64) *Bucket {
	b := &Bucket{rate: rate, capacity: capacity, tokens: make(chan struct{}, capacity)}
	go b.startRefill()

	return b
}

func (b *Bucket) startRefill() {
	for {
		select {
		case <-b.done:
			return

		default:
			if len(b.tokens) >= int(b.capacity) {
				break // break select, goes to sleep
			}

			b.tokens <- struct{}{}
		}

		time.Sleep(time.Second / time.Duration(b.rate))
	}
}

func (b *Bucket) Stop() {
	close(b.done)
	close(b.tokens)
}

func (b *Bucket) Take(n int) {
	for i := 0; i < n; i++ {
		<-b.tokens
	}
}

func (b *Bucket) TakeOne() {
	<-b.tokens
}

func (b *Bucket) Debug() {
	fmt.Printf("rate: %v, capacity: %v, num of tokens: %v\n", b.rate, b.capacity, len(b.tokens))
}
