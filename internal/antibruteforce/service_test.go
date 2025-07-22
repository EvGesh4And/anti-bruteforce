package antibruteforce

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/stretchr/testify/assert"
)

func newTestService() *Service {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.SecurityConfig{
		LoginRate:       2,
		PassRate:        2,
		IPRate:          2,
		CleanupInterval: time.Minute,
		BucketMaxIdle:   time.Minute,
	}
	return NewService(logger, cfg)
}

func TestService_CheckRateLimiting(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	ip := "192.168.1.10"
	login := "user1"
	pass := "pass1"

	// Allow 2 attempts
	assert.True(t, svc.Check(ctx, login, pass, ip))
	assert.True(t, svc.Check(ctx, login, pass, ip))

	// 3rd attempt should be blocked
	assert.False(t, svc.Check(ctx, login, pass, ip))
}

func TestService_Whitelist(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.AddToWhitelist("192.168.1.0/24")
	assert.NoError(t, err)

	assert.True(t, svc.Check(ctx, "login", "pass", "192.168.1.10")) // should always pass
}

func TestService_Blacklist(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.AddToBlacklist("10.0.0.0/8")
	assert.NoError(t, err)

	assert.False(t, svc.Check(ctx, "login", "pass", "10.1.1.1")) // should always fail
}

func TestService_Reset(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	ip := "172.16.0.1"
	login := "admin"

	// Первая и вторая попытки — валидные
	assert.True(t, svc.Check(ctx, login, "pass1", ip))
	assert.True(t, svc.Check(ctx, login, "pass2", ip))

	// 3-я попытка — превысит лимит login или ip
	assert.False(t, svc.Check(ctx, login, "pass3", ip))

	// Сброс login и ip, но НЕ password
	svc.Reset(login, ip)

	// Новая попытка с новым паролем (чтобы пароль не блокировал)
	assert.True(t, svc.Check(ctx, login, "pass4", ip))
}

func TestService_InvalidIP(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	ok := svc.Check(ctx, "login", "pass", "invalid_ip")
	assert.False(t, ok)
}
