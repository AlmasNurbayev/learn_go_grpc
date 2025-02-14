package authGrpc

import (
	"context"
	"errors"
	"log/slog"
	"sso/internal/errorsPackage"
	"sso/internal/services/auth"

	ssov1 "github.com/AlmasNurbayev/learn_go_grpc_protos/generated/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth auth.AuthService
	log  *slog.Logger
}

func Register(gRPC *grpc.Server, authService *auth.AuthService, log *slog.Logger) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: *authService, log: log})
}

func (s *serverAPI) Login(ctx context.Context, data *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	payload := LoginRequestForValidate{
		Login:    data.Login,
		Password: data.Password,
		Type:     data.Type,
		AppId:    data.AppId,
	}

	//fmt.Println("handler", ctx.Value(middleware.TraceIDKey))

	// TODO - условная валидация, если type = "email", то проверять что поле login является email
	err := ValidateStruct(payload)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, data.Login, data.Type, data.Password, int(data.AppId))
	if err != nil {
		if errors.Is(err, errorsPackage.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, errorsPackage.ErrInvalidCredentials.Error())
		}
		if errors.Is(err, errorsPackage.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, errorsPackage.ErrInvalidCredentials.Error())
		}
		if errors.Is(err, errorsPackage.ErrAppNotFound) {
			return nil, status.Error(codes.InvalidArgument, errorsPackage.ErrAppNotFound.Error())
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
		if errors.Is(err, errorsPackage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, errorsPackage.ErrUserExists.Error())
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

	isAdmin, err := s.auth.IsAdmin(ctx, data.UserId)
	if err != nil {
		if errors.Is(err, errorsPackage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errorsPackage.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil

}
