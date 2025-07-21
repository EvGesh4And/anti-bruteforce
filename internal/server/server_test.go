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
	CheckFn               func(ctx context.Context, login, password, ip string) bool
	ResetFn               func(login, ip string)
	AddToBlacklistFn      func(cidr string) error
	RemoveFromBlacklistFn func(cidr string)
	AddToWhitelistFn      func(cidr string) error
	RemoveFromWhitelistFn func(cidr string)
}

func (m *mockService) Check(ctx context.Context, login, password, ip string) bool {
	return m.CheckFn(ctx, login, password, ip)
}

func (m *mockService) Reset(login, ip string) {
	m.ResetFn(login, ip)
}

func (m *mockService) AddToBlacklist(cidr string) error {
	return m.AddToBlacklistFn(cidr)
}

func (m *mockService) RemoveFromBlacklist(cidr string) {
	m.RemoveFromBlacklistFn(cidr)
}

func (m *mockService) AddToWhitelist(cidr string) error {
	return m.AddToWhitelistFn(cidr)
}

func (m *mockService) RemoveFromWhitelist(cidr string) {
	m.RemoveFromWhitelistFn(cidr)
}

func newTestServer(mock *mockService) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return New(logger, mock)
}

func TestCheck(t *testing.T) {
	mock := &mockService{
		CheckFn: func(ctx context.Context, login, password, ip string) bool {
			_ = ctx
			assert.Equal(t, "user", login)
			assert.Equal(t, "pass", password)
			assert.Equal(t, "127.0.0.1", ip)
			return true
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

func TestReset(t *testing.T) {
	called := false
	mock := &mockService{
		ResetFn: func(login, ip string) {
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

func TestAddToBlacklist(t *testing.T) {
	mock := &mockService{
		AddToBlacklistFn: func(cidr string) error {
			assert.Equal(t, "192.168.0.0/24", cidr)
			return nil
		},
	}

	server := newTestServer(mock)

	_, err := server.AddToBlacklist(context.Background(), &pb.NetworkRequest{
		Network: "192.168.0.0/24",
	})

	assert.NoError(t, err)
}

func TestAddToBlacklist_Error(t *testing.T) {
	mock := &mockService{
		AddToBlacklistFn: func(cidr string) error {
			_ = cidr
			return errors.New("bad CIDR")
		},
	}

	server := newTestServer(mock)

	_, err := server.AddToBlacklist(context.Background(), &pb.NetworkRequest{
		Network: "bad",
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "bad CIDR")
}

func TestRemoveFromBlacklist(t *testing.T) {
	called := false
	mock := &mockService{
		RemoveFromBlacklistFn: func(cidr string) {
			called = true
			assert.Equal(t, "10.0.0.0/8", cidr)
		},
	}

	server := newTestServer(mock)

	_, err := server.RemoveFromBlacklist(context.Background(), &pb.NetworkRequest{
		Network: "10.0.0.0/8",
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAddToWhitelist(t *testing.T) {
	mock := &mockService{
		AddToWhitelistFn: func(cidr string) error {
			assert.Equal(t, "172.16.0.0/12", cidr)
			return nil
		},
	}

	server := newTestServer(mock)

	_, err := server.AddToWhitelist(context.Background(), &pb.NetworkRequest{
		Network: "172.16.0.0/12",
	})

	assert.NoError(t, err)
}

func TestAddToWhitelist_Error(t *testing.T) {
	mock := &mockService{
		AddToWhitelistFn: func(cidr string) error {
			_ = cidr
			return errors.New("invalid")
		},
	}

	server := newTestServer(mock)

	_, err := server.AddToWhitelist(context.Background(), &pb.NetworkRequest{
		Network: "bad",
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid")
}

func TestRemoveFromWhitelist(t *testing.T) {
	called := false
	mock := &mockService{
		RemoveFromWhitelistFn: func(cidr string) {
			called = true
			assert.Equal(t, "192.168.1.0/24", cidr)
		},
	}

	server := newTestServer(mock)

	_, err := server.RemoveFromWhitelist(context.Background(), &pb.NetworkRequest{
		Network: "192.168.1.0/24",
	})

	assert.NoError(t, err)
	assert.True(t, called)
}
