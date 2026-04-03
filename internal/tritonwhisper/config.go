package tritonwhisper

import (
	"os"
	"strings"
)

// Config задаёт имена моделей в Triton и сетевые параметры клиента.
// Значения по умолчанию совпадают с deployments/model_repo.
type Config struct {
	GRPCAddress string

	EnsembleModel       string
	DecoderModel        string
	DecoderWithPastModel string
	ModelVersion        string

	ConnectTimeoutSec float64
	NetworkTimeoutSec float64
}

func (c *Config) setDefaults() {
	if c.GRPCAddress == "" {
		c.GRPCAddress = os.Getenv("TRITON_GRPC_ADDR")
	}
	if c.GRPCAddress == "" {
		c.GRPCAddress = "127.0.0.1:8001"
	}
	c.GRPCAddress = strings.TrimPrefix(c.GRPCAddress, "http://")
	c.GRPCAddress = strings.TrimPrefix(c.GRPCAddress, "https://")

	if c.EnsembleModel == "" {
		c.EnsembleModel = "whisper_ensemble"
	}
	if c.DecoderModel == "" {
		c.DecoderModel = "whisper_medium_fp16_decoder"
	}
	if c.DecoderWithPastModel == "" {
		c.DecoderWithPastModel = "whisper_medium_fp16_decoder_with_past"
	}
	if c.ConnectTimeoutSec <= 0 {
		c.ConnectTimeoutSec = 30
	}
	if c.NetworkTimeoutSec <= 0 {
		c.NetworkTimeoutSec = 120
	}
}
