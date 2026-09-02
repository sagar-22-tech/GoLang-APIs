package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBucketFull = errors.New("leaky bucket is full")
	ErrCancelled  = errors.New("request cancelled")
)

type RequestResult int

const (
	Released RequestResult = iota
	Expired
)

type Request struct {
	id        string
	createdAt time.Time
	done      chan struct{}
	result    RequestResult
}
type LeakyBucket struct {
	capacity int
	rate     time.Duration
	expiry   time.Duration

	queue  []*Request
	stopCh chan struct{}

	mu   sync.Mutex
	once sync.Once
}

func NewLeakyBucket(capacity int, rate time.Duration, expiry time.Duration) *LeakyBucket {
	var lb LeakyBucket
	lb.capacity = capacity
	lb.rate = rate
	lb.expiry = expiry
	lb.stopCh = make(chan struct{})
	return &lb

}

func (lb *LeakyBucket) Allow(ctx context.Context) (*Request, error) {
	lb.mu.Lock()

	if len(lb.queue) == lb.capacity {
		lb.mu.Unlock()
		return nil, ErrBucketFull
	}

	req := &Request{
		id:        uuid.New().String(),
		createdAt: time.Now(),
		done:      make(chan struct{}),
	}

	lb.queue = append(lb.queue, req)

	lb.mu.Unlock()

	select {
	case <-req.done:
		return req, nil

	case <-ctx.Done():
		lb.mu.Lock()
		lb.removeRequest(req)
		lb.mu.Unlock()

		return nil, ErrCancelled
	}
}
func (lb *LeakyBucket) removeRequest(req *Request) {
	for i, r := range lb.queue {
		if r == req {
			lb.queue = append(
				lb.queue[:i],
				lb.queue[i+1:]...,
			)
			return
		}
	}
}
func (lb *LeakyBucket) Leak() {
	ticker := time.NewTicker(lb.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lb.mu.Lock()

			// Remove expired requests
			for len(lb.queue) > 0 &&
				time.Since(lb.queue[0].createdAt) >= lb.expiry {

				req := lb.queue[0]

				req.result = Expired

				fmt.Println("expired:", req.id)

				lb.queue = lb.queue[1:]

				close(req.done)
			}

			// Leak one valid request
			if len(lb.queue) > 0 {
				req := lb.queue[0]

				req.result = Released

				fmt.Println("leaked:", req.id)

				lb.queue = lb.queue[1:]

				close(req.done)
			}

			lb.mu.Unlock()

		case <-lb.stopCh:
			return
		}
	}
}
func (lb *LeakyBucket) Stop() {
	lb.once.Do(func() {
		close(lb.stopCh)
	})
}
func RateLimitMiddleware(
	bucket *LeakyBucket,
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req, err := bucket.Allow(r.Context())

		if err != nil {
			if errors.Is(err, ErrBucketFull) {
				http.Error(
					w,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)
				return
			}

			if errors.Is(err, ErrCancelled) {
				return
			}

			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		// Request waited successfully, but check if it expired.
		if req.result == Expired {
			http.Error(
				w,
				"request expired",
				http.StatusTooManyRequests,
			)
			return
		}

		// Request was released by the bucket.
		next.ServeHTTP(w, r)
	})
}
