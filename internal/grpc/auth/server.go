package auth

import (
	"context"
	"errors"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/services/auth"
	authv1 "github.com/EtoNeAnanasbI95/protos_auth/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, login string, password string) (int64, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Validate(ctx context.Context, token string) (int64, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth AuthService
}

func Register(gRPC *grpc.Server, auth AuthService) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, login *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if login.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login required")
	}
	if login.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	refreshToken, accessToken, err := s.auth.Login(ctx, login.GetLogin(), login.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUserCredentials) {
			return nil, status.Error(codes.Unauthenticated, auth.ErrInvalidUserCredentials.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, register *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if register.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login required")
	}
	if register.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}
	uid, err := s.auth.Register(ctx, register.GetLogin(), register.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, auth.ErrUserAlreadyExists.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, token *authv1.TokenRequest) (*authv1.LoginResponse, error) {
	if token.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}

	refreshToken, accessToken, err := s.auth.Refresh(ctx, token.GetToken())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUserCredentials) {
			return nil, status.Error(codes.Unauthenticated, auth.ErrInvalidUserCredentials.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *serverAPI) Validate(ctx context.Context, token *authv1.TokenRequest) (*emptypb.Empty, error) {
	if token.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}

	_, err := s.auth.Validate(ctx, token.GetToken())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUserCredentials) {
			return nil, status.Error(codes.Unauthenticated, auth.ErrInvalidUserCredentials.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
