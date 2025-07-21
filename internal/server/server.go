package grpcserver

import (
	"context"
	"log/slog"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AntiBruteforceService interface {
	Check(ctx context.Context, login, password, ip string) bool
	Reset(login, ip string)
	AddToBlacklist(network string) error
	RemoveFromBlacklist(network string)
	AddToWhitelist(network string) error
	RemoveFromWhitelist(network string)
}

type Server struct {
	pb.UnimplementedAntiBruteforceServer
	service AntiBruteforceService // Используем интерфейс
	logger  *slog.Logger
}

func New(logger *slog.Logger, service AntiBruteforceService) *Server {
	return &Server{logger: logger, service: service}
}

func (s *Server) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	ok := s.service.Check(ctx, req.Login, req.Password, req.Ip)
	return &pb.CheckResponse{Ok: ok}, nil
}

func (s *Server) Reset(ctx context.Context, req *pb.ResetRequest) (*emptypb.Empty, error) {
	_ = ctx
	s.service.Reset(req.Login, req.Ip)
	return &emptypb.Empty{}, nil
}

func (s *Server) AddToBlacklist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	_ = ctx
	if err := s.service.AddToBlacklist(req.Network); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveFromBlacklist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	_ = ctx
	s.service.RemoveFromBlacklist(req.Network)
	return &emptypb.Empty{}, nil
}

func (s *Server) AddToWhitelist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	_ = ctx
	if err := s.service.AddToWhitelist(req.Network); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveFromWhitelist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	_ = ctx
	s.service.RemoveFromWhitelist(req.Network)
	return &emptypb.Empty{}, nil
}
