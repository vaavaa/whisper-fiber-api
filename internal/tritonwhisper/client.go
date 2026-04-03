package tritonwhisper

import (
	"fmt"

	"github.com/Trendyol/go-triton-client/base"
	tritongrpc "github.com/Trendyol/go-triton-client/client/grpc"
)

// Client оборачивает gRPC Triton и конфиг имён моделей.
type Client struct {
	cfg    Config
	triton base.Client
}

// NewClient создаёт клиент (без TLS, как в локальном docker-compose).
func NewClient(cfg Config) (*Client, error) {
	cfg.setDefaults()
	tc, err := tritongrpc.NewClient(
		cfg.GRPCAddress,
		false,
		cfg.ConnectTimeoutSec,
		cfg.NetworkTimeoutSec,
		false,
		false,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("triton grpc: %w", err)
	}
	return &Client{cfg: cfg, triton: tc}, nil
}
