package bucket

import (
	"sync"
	"testing"
	"time"
)

func TestAllowBasic(t *testing.T) {
	m := NewManager(3, time.Second) // 3 req/sec
	key := "client"

	for i := 0; i < 3; i++ {
		if !m.Allow(key) {
			t.Fatalf("request #%d should be allowed", i+1)
		}
	}
	if m.Allow(key) {
		t.Fatal("4th request should be denied due to rate limit")
	}
}

func TestAllowLeakOverTime(t *testing.T) {
	m := NewManager(2, time.Second) // 2 req/sec
	key := "slow-client"

	m.Allow(key)
	m.Allow(key)
	time.Sleep(600 * time.Millisecond)

	if !m.Allow(key) {
		t.Fatal("should allow after partial leak")
	}

	time.Sleep(1 * time.Second)

	if !m.Allow(key) {
		t.Fatal("should allow after full recovery")
	}
}

func TestReset(t *testing.T) {
	m := NewManager(1, time.Second)
	key := "reset-me"

	if !m.Allow(key) {
		t.Fatal("first request should be allowed")
	}
	if m.Allow(key) {
		t.Fatal("second request should be denied")
	}

	m.Reset(key)

	if !m.Allow(key) {
		t.Fatal("should be allowed after reset")
	}
}

func TestSeparateKeys(t *testing.T) {
	m := NewManager(1, time.Second)

	if !m.Allow("user1") {
		t.Fatal("user1 first request should be allowed")
	}
	if !m.Allow("user2") {
		t.Fatal("user2 first request should be allowed")
	}
	if m.Allow("user1") {
		t.Fatal("user1 second request should be denied")
	}
}

func TestParallelAccess(t *testing.T) {
	m := NewManager(100, time.Second)
	key := "parallel-client"

	var wg sync.WaitGroup
	success := 0
	mu := sync.Mutex{}

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.Allow(key) {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if success > 100 {
		t.Fatalf("too many requests allowed: %d > 100", success)
	}
	if success < 90 {
		t.Fatalf("too few requests allowed (expected ~100): %d", success)
	}
}
