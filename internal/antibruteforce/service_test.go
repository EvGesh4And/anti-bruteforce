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

	ok, err := svc.Check(ctx, login, pass, ip)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = svc.Check(ctx, login, pass, ip)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = svc.Check(ctx, login, pass, ip)
	assert.False(t, ok)
	assert.NoError(t, err)
}

func TestService_Whitelist(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.AddToWhitelist(ctx, "192.168.1.0/24")
	assert.NoError(t, err)

	ok, err := svc.Check(ctx, "login", "pass", "192.168.1.10")
	assert.True(t, ok)
	assert.NoError(t, err)

	err = svc.RemoveFromWhitelist(ctx, "192.168.1.0/24")
	assert.NoError(t, err)

	// После удаления из белого списка — ограничение по лимитам
	ok, err = svc.Check(ctx, "login", "pass", "192.168.1.10")
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestService_Blacklist(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.AddToBlacklist(ctx, "10.0.0.0/8")
	assert.NoError(t, err)

	ok, err := svc.Check(ctx, "login", "pass", "10.1.1.1")
	assert.False(t, ok)
	assert.NoError(t, err)

	err = svc.RemoveFromBlacklist(ctx, "10.0.0.0/8")
	assert.NoError(t, err)

	ok, err = svc.Check(ctx, "login", "pass", "10.1.1.1")
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestService_Reset(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	ip := "172.16.0.1"
	login := "admin"

	ok, err := svc.Check(ctx, login, "pass1", ip)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = svc.Check(ctx, login, "pass2", ip)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = svc.Check(ctx, login, "pass3", ip)
	assert.False(t, ok)
	assert.NoError(t, err)

	svc.Reset(ctx, login, ip)

	ok, err = svc.Check(ctx, login, "pass4", ip)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestService_InvalidIP(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	ok, err := svc.Check(ctx, "login", "pass", "invalid_ip")
	assert.False(t, ok)
	assert.ErrorIs(t, err, ErrInvalidIP)
}

func TestService_InvalidCIDR(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.AddToWhitelist(ctx, "invalid_cidr")
	assert.ErrorIs(t, err, ErrInvalidNetwork)

	err = svc.AddToBlacklist(ctx, "invalid_cidr")
	assert.ErrorIs(t, err, ErrInvalidNetwork)

	err = svc.RemoveFromWhitelist(ctx, "invalid_cidr")
	assert.ErrorIs(t, err, ErrInvalidNetwork)

	err = svc.RemoveFromBlacklist(ctx, "invalid_cidr")
	assert.ErrorIs(t, err, ErrInvalidNetwork)
}

func TestService_Shutdown(_ *testing.T) {
	svc := newTestService()
	svc.Shutdown()
	// Здесь просто проверяем, что метод выполняется без паники
}
