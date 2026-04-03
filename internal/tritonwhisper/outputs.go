package tritonwhisper

import (
	"fmt"

	"github.com/Trendyol/go-triton-client/base"
	tritongrpc "github.com/Trendyol/go-triton-client/client/grpc"
)

func outAllOutputs() []base.InferOutput {
	names := []string{"logits"}
	for i := range DecoderLayers {
		p := fmt.Sprintf("present.%d", i)
		names = append(names,
			p+".decoder.key", p+".decoder.value",
			p+".encoder.key", p+".encoder.value",
		)
	}
	req := make([]base.InferOutput, 0, len(names))
	for _, n := range names {
		req = append(req, tritongrpc.NewInferOutput(n, map[string]any{"binary_data": true}))
	}
	return req
}
