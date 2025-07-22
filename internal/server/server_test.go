package grpcserver

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	CheckFn               func(ctx context.Context, login, password, ip string) (bool, error)
	ResetFn               func(ctx context.Context, login, ip string)
	AddToBlacklistFn      func(ctx context.Context, cidr string) error
	RemoveFromBlacklistFn func(ctx context.Context, cidr string) error
	AddToWhitelistFn      func(ctx context.Context, cidr string) error
	RemoveFromWhitelistFn func(ctx context.Context, cidr string) error
}

func (m *mockService) Check(ctx context.Context, login, password, ip string) (bool, error) {
	return m.CheckFn(ctx, login, password, ip)
}

func (m *mockService) Reset(ctx context.Context, login, ip string) {
	m.ResetFn(ctx, login, ip)
}

func (m *mockService) AddToBlacklist(ctx context.Context, cidr string) error {
	return m.AddToBlacklistFn(ctx, cidr)
}

func (m *mockService) RemoveFromBlacklist(ctx context.Context, cidr string) error {
	return m.RemoveFromBlacklistFn(ctx, cidr)
}

func (m *mockService) AddToWhitelist(ctx context.Context, cidr string) error {
	return m.AddToWhitelistFn(ctx, cidr)
}

func (m *mockService) RemoveFromWhitelist(ctx context.Context, cidr string) error {
	return m.RemoveFromWhitelistFn(ctx, cidr)
}

func newTestServer(mock *mockService) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return New(logger, mock)
}

func TestCheck(t *testing.T) {
	mock := &mockService{
		CheckFn: func(_ context.Context, login, password, ip string) (bool, error) {
			assert.Equal(t, "user", login)
			assert.Equal(t, "pass", password)
			assert.Equal(t, "127.0.0.1", ip)
			return true, nil
		},
	}

	server := newTestServer(mock)

	resp, err := server.Check(context.Background(), &pb.CheckRequest{
		Login:    "user",
		Password: "pass",
		Ip:       "127.0.0.1",
	})

	assert.NoError(t, err)
	assert.True(t, resp.Ok)
}

func TestCheck_Limited(t *testing.T) {
	mock := &mockService{
		CheckFn: func(_ context.Context, _, _, _ string) (bool, error) {
			return false, nil // не ошибка, просто отказ
		},
	}

	server := newTestServer(mock)

	resp, err := server.Check(context.Background(), &pb.CheckRequest{
		Login:    "user",
		Password: "pass",
		Ip:       "127.0.0.1",
	})

	assert.NoError(t, err)
	assert.False(t, resp.Ok)
}

func TestCheck_InvalidIP(t *testing.T) {
	mock := &mockService{
		CheckFn: func(_ context.Context, _, _, _ string) (bool, error) {
			return false, errors.New("invalid IP address")
		},
	}

	server := newTestServer(mock)

	resp, err := server.Check(context.Background(), &pb.CheckRequest{
		Login:    "user",
		Password: "pass",
		Ip:       "invalid_ip",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid IP address")
}

func TestReset(t *testing.T) {
	called := false
	mock := &mockService{
		ResetFn: func(_ context.Context, login, ip string) {
			called = true
			assert.Equal(t, "user", login)
			assert.Equal(t, "127.0.0.1", ip)
		},
	}

	server := newTestServer(mock)

	_, err := server.Reset(context.Background(), &pb.ResetRequest{
		Login: "user",
		Ip:    "127.0.0.1",
	})

	assert.NoError(t, err)
	assert.True(t, called)
}
