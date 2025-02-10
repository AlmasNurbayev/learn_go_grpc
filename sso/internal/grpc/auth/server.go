package authGrpc

import (
	"context"
	"errors"
	"sso/internal/services/auth"
	"sso/internal/storage"

	ssov1 "github.com/AlmasNurbayev/learn_go_grpc_protos/generated/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth auth.AuthService
}

func Register(gRPC *grpc.Server, authService *auth.AuthService) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: *authService})
}

func (s *serverAPI) Login(ctx context.Context, data *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	payload := LoginRequestForValidate{
		Login:    data.Login,
		Password: data.Password,
		Type:     data.Type,
		AppId:    data.AppId,
	}
	// TODO - условная валидация, если type = "email", то проверять что поле login является email
	err := ValidateStruct(payload)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, data.Login, data.Type, data.Password, int(data.AppId))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		if errors.Is(err, auth.ErrInvalidAppId) {
			return nil, status.Error(codes.InvalidArgument, "invalid app id")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, data *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	payload := RegisterRequestForValidate{
		Password: data.Password,
		Email:    data.Email,
		Phone:    data.Phone,
	}
	err := ValidateStruct(payload)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userId, err := s.auth.RegisterNewUser(ctx, data.Email, data.Phone, data.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, data *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	payload := IsAdminForValidate{
		UserId: data.UserId,
	}
	err := ValidateStruct(payload)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if data.UserId == 10 {
		return &ssov1.IsAdminResponse{IsAdmin: true}, nil
	} else {
		return &ssov1.IsAdminResponse{IsAdmin: false}, nil
	}
}
