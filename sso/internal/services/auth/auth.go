package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger"
	"sso/internal/models"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log         *slog.Logger
	storage     AuthStorage
	tokenTTL    time.Duration
	appProvider AppProvider
}

type AuthStorage interface {
	SaveUser(ctx context.Context, email string, phone string, passHash []byte) (id int64, err error)
	GetUserByPhone(ctx context.Context, phone string) (user models.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
	GetUserById(ctx context.Context, id int64) (user models.User, err error)
	IsAdmin(ctx context.Context, id int64) (isAdmin bool, err error)
	GetAppById(ctx context.Context, id int) (app models.App, err error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
)

func NewService(log *slog.Logger, storage AuthStorage, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, login string,
	typeLogin string, password string, appId int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))

	log.Info("login user ", slog.String("login", login))

	var user models.User
	var err error

	if typeLogin == "email" {
		user, err = a.storage.GetUserByEmail(ctx, login)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Warn("user not found", slog.String("email", login))
				return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
			}
			log.Error("failed to get user by email", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}
		log.Info("login user by email", slog.Int64("id", user.Id))
	}

	if typeLogin == "phone" {
		user, err = a.storage.GetUserByPhone(ctx, login)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Warn("user not found", slog.String("email", login))
				return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
			}
			log.Error("failed to get user by phone", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}
		log.Info("login user by email", slog.Int64("id", user.Id))
	}

	app, err := a.storage.GetAppById(ctx, int(appId))
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", slog.Int("id", appId))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}
		log.Error("failed to get user by phone", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully", slog.Int64("id", user.Id))

	token, err := jwt.GenerateToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("token generated for user.id", slog.Int64("id", user.Id))

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email string, phone string,
	password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))
	log.Info("registering new user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error())})
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := a.storage.SaveUser(ctx, email, phone, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("email", email))
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		log.Error("failed to save user", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("register new user ", slog.Int64("id", id))
	return id, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op))

	isAdmin, err := a.storage.IsAdmin(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.Int64("id", id))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		} else if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("failed to get app by id", logger.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		} else {
			log.Error("failed to get user by id", logger.Err(err))
			return false, fmt.Errorf("%s: %w", op, err)

		}
	}
	log.Info("user checked is admin ", slog.Int64("id ", id), slog.Bool("isAdmin ", isAdmin))
	return isAdmin, nil
}
