package bucket

import (
	"context"
	"sync"
	"time"
)

// leakyBucket enforces rate limit using the leaky bucket algorithm.
type leakyBucket struct {
	mu       sync.Mutex
	last     time.Time
	tokens   float64
	capacity float64
	leakRate float64 // tokens per second
}

func newLeakyBucket(limit int, interval time.Duration) *leakyBucket {
	return &leakyBucket{
		last:     time.Now(),
		capacity: float64(limit),
		leakRate: float64(limit) / interval.Seconds(),
	}
}

func (b *leakyBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	// Leak tokens over time
	b.tokens -= b.leakRate * elapsed
	if b.tokens < 0 {
		b.tokens = 0
	}

	if b.tokens+1 > b.capacity {
		return false
	}

	b.tokens++
	return true
}

// Manager stores leaky buckets by key.
type Manager struct {
	mu      sync.Mutex
	buckets map[string]*leakyBucket
	limit   int
	window  time.Duration
}

func NewManager(limit int, window time.Duration) *Manager {
	return &Manager{
		buckets: make(map[string]*leakyBucket),
		limit:   limit,
		window:  window,
	}
}

func (m *Manager) Allow(key string) bool {
	m.mu.Lock()
	b, ok := m.buckets[key]
	if !ok {
		b = newLeakyBucket(m.limit, m.window)
		m.buckets[key] = b
	}
	m.mu.Unlock()

	return b.allow()
}

func (m *Manager) Reset(key string) {
	m.mu.Lock()
	delete(m.buckets, key)
	m.mu.Unlock()
}

func (m *Manager) StartCleanup(ctx context.Context, interval, maxIdle time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				m.mu.Lock()
				for key, b := range m.buckets {
					select {
					case <-ctx.Done():
						defer m.mu.Unlock()
						return
					default:
						b.mu.Lock()
						if now.Sub(b.last) > maxIdle {
							delete(m.buckets, key)
						}
						b.mu.Unlock()
					}
				}
				m.mu.Unlock()
			}
		}
	}()
}
