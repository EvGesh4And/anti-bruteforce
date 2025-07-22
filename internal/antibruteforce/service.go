package antibruteforce

import (
	"context"
	"log/slog"
	"net/netip"
	"sync"
	"time"

	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/EvGesh4And/anti-bruteforce/internal/bucket"
)

type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
}

// Service implements anti-bruteforce logic.
type Service struct {
	logger       *slog.Logger
	loginBuckets RateLimiter
	passBuckets  RateLimiter
	ipBuckets    RateLimiter

	mu         sync.RWMutex
	whitelist  map[string]netip.Prefix
	blacklist  map[string]netip.Prefix
	cancelFunc context.CancelFunc
}

// NewService creates a new Service with given limits per minute.
func NewService(logger *slog.Logger, cfg config.SecurityConfig) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	login := bucket.NewManager(cfg.LoginRate, time.Minute)
	pass := bucket.NewManager(cfg.PassRate, time.Minute)
	ip := bucket.NewManager(cfg.IPRate, time.Minute)

	login.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)
	pass.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)
	ip.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)

	return &Service{
		logger:       logger,
		loginBuckets: login,
		passBuckets:  pass,
		ipBuckets:    ip,
		whitelist:    make(map[string]netip.Prefix),
		blacklist:    make(map[string]netip.Prefix),
		cancelFunc:   cancel,
	}
}

// Check authorisation attempt.
func (s *Service) Check(ctx context.Context, login, password, ipStr string) bool {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		s.logger.ErrorContext(ctx, "invalid ip", "ip", ipStr)
		return false
	}

	if s.isWhitelisted(ip) {
		return true
	}
	if s.isBlacklisted(ip) {
		return false
	}

	if !s.loginBuckets.Allow(login) {
		s.logger.WarnContext(ctx, "login rate limit exceeded", "login", login)
		return false
	}
	if !s.passBuckets.Allow(password) {
		s.logger.WarnContext(ctx, "password rate limit exceeded")
		return false
	}
	if !s.ipBuckets.Allow(ipStr) {
		s.logger.WarnContext(ctx, "ip rate limit exceeded", "ip", ipStr)
		return false
	}
	return true
}

// Reset cleans buckets for login and ip.
func (s *Service) Reset(login, ip string) {
	s.loginBuckets.Reset(login)
	s.ipBuckets.Reset(ip)
}

func (s *Service) isWhitelisted(ip netip.Addr) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.whitelist {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

func (s *Service) isBlacklisted(ip netip.Addr) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.blacklist {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

func (s *Service) addToList(list map[string]netip.Prefix, network string) error {
	p, err := netip.ParsePrefix(network)
	if err != nil {
		return err
	}
	s.mu.Lock()
	list[p.String()] = p
	s.mu.Unlock()
	return nil
}

func (s *Service) removeFromList(list map[string]netip.Prefix, network string) {
	p, err := netip.ParsePrefix(network)
	if err != nil {
		return
	}
	s.mu.Lock()
	delete(list, p.String())
	s.mu.Unlock()
}

// AddToBlacklist appends a network to the blacklist.
func (s *Service) AddToBlacklist(network string) error {
	return s.addToList(s.blacklist, network)
}

// RemoveFromBlacklist removes a network from the blacklist.
func (s *Service) RemoveFromBlacklist(network string) {
	s.removeFromList(s.blacklist, network)
}

// AddToWhitelist appends a network to the whitelist.
func (s *Service) AddToWhitelist(network string) error {
	return s.addToList(s.whitelist, network)
}

// RemoveFromWhitelist removes a network from the whitelist.
func (s *Service) RemoveFromWhitelist(network string) {
	s.removeFromList(s.whitelist, network)
}

func (s *Service) Shutdown() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
