package tritonwhisper

import (
	"context"
	"fmt"

	"github.com/Trendyol/go-triton-client/base"
	tritongrpc "github.com/Trendyol/go-triton-client/client/grpc"
	"github.com/Trendyol/go-triton-client/options"
)

// EncoderOut — плоский тензор last_hidden_state энкодера и его форма [1,1500,1024].
type EncoderOut struct {
	Values []float32
	Shape  []int64
}

// RunWhisperEnsemble: AUDIO_PCM (float32 mono 16 kHz) → encoder_hidden_states.
func (c *Client) RunWhisperEnsemble(ctx context.Context, pcmFloat32 []float32) (*EncoderOut, error) {
	if len(pcmFloat32) == 0 {
		return nil, fmt.Errorf("пустой PCM")
	}
	in := tritongrpc.NewInferInput("AUDIO_PCM", "FP32", []int64{int64(len(pcmFloat32))}, nil)
	if err := in.SetData(pcmFloat32, true); err != nil {
		return nil, fmt.Errorf("audio tensor: %w", err)
	}
	outputs := []base.InferOutput{
		tritongrpc.NewInferOutput("encoder_hidden_states", map[string]any{"binary_data": true}),
	}
	res, err := c.triton.Infer(ctx, c.cfg.EnsembleModel, c.cfg.ModelVersion, []base.InferInput{in}, outputs, &options.InferOptions{})
	if err != nil {
		return nil, err
	}
	shape, err := res.GetShape("encoder_hidden_states")
	if err != nil {
		return nil, err
	}
	flat, err := res.AsFloat32Slice("encoder_hidden_states")
	if err != nil {
		return nil, err
	}
	return &EncoderOut{Values: flat, Shape: cloneShape(shape)}, nil
}

func cloneShape(s []int64) []int64 {
	out := make([]int64, len(s))
	copy(out, s)
	return out
}

// ArgmaxLastTimeStep — для logits формы [1, S, V] берёт argmax по последнему S.
func ArgmaxLastTimeStep(logits []float32, shape []int64) (int64, error) {
	if len(shape) != 3 || shape[0] != 1 {
		return 0, fmt.Errorf("logits: ожидалась форма [1,S,V], получено %v", shape)
	}
	s := shape[1]
	v := shape[2]
	if int(s*v) != len(logits) {
		return 0, fmt.Errorf("logits: длина %d не совпадает с shape %v", len(logits), shape)
	}
	start := int((s - 1) * v)
	row := logits[start:]
	best := 0
	for i := 1; i < len(row); i++ {
		if row[i] > row[best] {
			best = i
		}
	}
	return int64(best), nil
}
