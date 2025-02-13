package suite

import (
	"context"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/AlmasNurbayev/learn_go_grpc_protos/generated/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const grpcHost = "localhost"

func New(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	cfg := config.ReadConfigByPath("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.NewClient(grpcAdress(&cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create grpc client: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        &cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAdress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
