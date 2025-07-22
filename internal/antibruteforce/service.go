package antibruteforce

import (
	"context"
	"errors"
	"log/slog"
	"net/netip"
	"sync"
	"time"

	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/EvGesh4And/anti-bruteforce/internal/bucket"
	"github.com/EvGesh4And/anti-bruteforce/internal/logger"
)

var (
	ErrInvalidIP      = errors.New("invalid IP address")
	ErrInvalidNetwork = errors.New("invalid CIDR network")
)

type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
}

type Service struct {
	lg           *slog.Logger
	loginBuckets RateLimiter
	passBuckets  RateLimiter
	ipBuckets    RateLimiter

	mu         sync.RWMutex
	whitelist  map[string]netip.Prefix
	blacklist  map[string]netip.Prefix
	cancelFunc context.CancelFunc
}

func NewService(lg *slog.Logger, cfg config.SecurityConfig) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	login := bucket.NewManager(cfg.LoginRate, time.Minute)
	pass := bucket.NewManager(cfg.PassRate, time.Minute)
	ip := bucket.NewManager(cfg.IPRate, time.Minute)

	login.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)
	pass.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)
	ip.StartCleanup(ctx, cfg.CleanupInterval, cfg.BucketMaxIdle)

	return &Service{
		lg:           lg,
		loginBuckets: login,
		passBuckets:  pass,
		ipBuckets:    ip,
		whitelist:    make(map[string]netip.Prefix),
		blacklist:    make(map[string]netip.Prefix),
		cancelFunc:   cancel,
	}
}

func (s *Service) setLogCompMeth(ctx context.Context, method string) context.Context {
	ctx = logger.WithLogComponent(ctx, "service")
	return logger.WithLogMethod(ctx, method)
}

func (s *Service) Check(ctx context.Context, login, password, ipStr string) (bool, error) {
	ctx = s.setLogCompMeth(ctx, "Check")
	ctx = logger.WithLogIP(ctx, ipStr)
	ctx = logger.WithLogLogin(ctx, login)

	s.lg.DebugContext(ctx, "starting authorization check")

	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, logger.AddPrefix(ctx, ErrInvalidIP)
	}

	if s.isWhitelisted(ip) {
		s.lg.InfoContext(ctx, "IP is whitelisted, authorization allowed")
		return true, nil
	}

	if s.isBlacklisted(ip) {
		s.lg.InfoContext(ctx, "IP is blacklisted, access denied")
		return false, nil
	}

	if !s.loginBuckets.Allow(login) {
		s.lg.InfoContext(ctx, "login rate limit exceeded")
		return false, nil
	}

	if !s.passBuckets.Allow(password) {
		s.lg.InfoContext(ctx, "password rate limit exceeded")
		return false, nil
	}

	if !s.ipBuckets.Allow(ipStr) {
		s.lg.InfoContext(ctx, "IP rate limit exceeded")
		return false, nil
	}

	s.lg.InfoContext(ctx, "authorization allowed")
	return true, nil
}

func (s *Service) Reset(ctx context.Context, login, ip string) {
	ctx = s.setLogCompMeth(ctx, "Reset")
	s.lg.DebugContext(ctx, "resetting rate limits")
	s.loginBuckets.Reset(login)
	s.ipBuckets.Reset(ip)
	s.lg.InfoContext(ctx, "rate limits reset")
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

func (s *Service) parsePrefix(ctx context.Context, network string) (netip.Prefix, error) {
	ctx = s.setLogCompMeth(ctx, "parsePrefix")
	p, err := netip.ParsePrefix(network)
	if err != nil {
		return netip.Prefix{}, logger.AddPrefix(ctx, ErrInvalidNetwork)
	}
	return p, nil
}

func (s *Service) addToList(ctx context.Context, list map[string]netip.Prefix, network string) error {
	ctx = s.setLogCompMeth(ctx, "addToList")
	ctx = logger.WithLogNetwork(ctx, network)
	s.lg.DebugContext(ctx, "adding network to list")
	p, err := s.parsePrefix(ctx, network)
	if err != nil {
		return err
	}
	s.mu.Lock()
	list[p.String()] = p
	s.mu.Unlock()
	s.lg.InfoContext(ctx, "network added")
	return nil
}

func (s *Service) removeFromList(ctx context.Context, list map[string]netip.Prefix, network string) error {
	ctx = s.setLogCompMeth(ctx, "removeFromList")
	ctx = logger.WithLogNetwork(ctx, network)

	s.lg.DebugContext(ctx, "removing network from list")
	p, err := s.parsePrefix(ctx, network)
	if err != nil {
		return err
	}
	s.mu.Lock()
	delete(list, p.String())
	s.mu.Unlock()
	s.lg.InfoContext(ctx, "network removed")
	return nil
}

func (s *Service) AddToBlacklist(ctx context.Context, network string) error {
	ctx = s.setLogCompMeth(ctx, "AddToBlacklist")
	if err := s.addToList(ctx, s.blacklist, network); err != nil {
		return logger.AddPrefix(ctx, err)
	}
	return nil
}

func (s *Service) RemoveFromBlacklist(ctx context.Context, network string) error {
	ctx = s.setLogCompMeth(ctx, "RemoveFromBlacklist")
	if err := s.removeFromList(ctx, s.blacklist, network); err != nil {
		return logger.AddPrefix(ctx, err)
	}
	return nil
}

func (s *Service) AddToWhitelist(ctx context.Context, network string) error {
	ctx = s.setLogCompMeth(ctx, "AddToWhitelist")
	if err := s.addToList(ctx, s.whitelist, network); err != nil {
		return logger.AddPrefix(ctx, err)
	}
	return nil
}

func (s *Service) RemoveFromWhitelist(ctx context.Context, network string) error {
	ctx = s.setLogCompMeth(ctx, "RemoveFromWhitelist")
	if err := s.removeFromList(ctx, s.whitelist, network); err != nil {
		return logger.AddPrefix(ctx, err)
	}
	return nil
}

func (s *Service) Shutdown() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
