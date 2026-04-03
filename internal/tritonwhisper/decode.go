package tritonwhisper

import (
	"context"
	"fmt"

	"github.com/Trendyol/go-triton-client/base"
	"github.com/Trendyol/go-triton-client/options"
)

// GreedyDecodeTokenIDs — жадный цикл: decoder → decoder_with_past.
// prefix обычно DefaultDecoderPrefix; для другого языка подставьте свои forced tokens.
func (c *Client) GreedyDecodeTokenIDs(ctx context.Context, enc *EncoderOut, prefix []int64) ([]int64, error) {
	if len(prefix) == 0 {
		prefix = DefaultDecoderPrefix
	}
	res, err := c.decoderFirst(ctx, enc, prefix)
	if err != nil {
		return nil, err
	}
	shape, err := res.GetShape("logits")
	if err != nil {
		return nil, err
	}
	logits, err := res.AsFloat32Slice("logits")
	if err != nil {
		return nil, err
	}
	tok, err := ArgmaxLastTimeStep(logits, shape)
	if err != nil {
		return nil, err
	}

	out := make([]int64, 0, len(prefix)+MaxDecoderPositions)
	out = append(out, prefix...)
	out = append(out, tok)

	for tok != EOSTokenID && len(out) < MaxDecoderPositions {
		res, err = c.decoderWithPast(ctx, tok, res)
		if err != nil {
			return nil, err
		}
		shape, err = res.GetShape("logits")
		if err != nil {
			return nil, err
		}
		logits, err = res.AsFloat32Slice("logits")
		if err != nil {
			return nil, err
		}
		tok, err = ArgmaxLastTimeStep(logits, shape)
		if err != nil {
			return nil, err
		}
		out = append(out, tok)
	}
	return out, nil
}

// TranscribePCMToTokenIDs — ensemble + жадный декодер (без BPE → текст; см. комментарий к пакету).
func (c *Client) TranscribePCMToTokenIDs(ctx context.Context, pcm []float32) ([]int64, error) {
	enc, err := c.RunWhisperEnsemble(ctx, pcm)
	if err != nil {
		return nil, err
	}
	return c.GreedyDecodeTokenIDs(ctx, enc, nil)
}

func (c *Client) decoderFirst(ctx context.Context, enc *EncoderOut, prefix []int64) (base.InferResult, error) {
	idsIn, err := int64TensorInput("input_ids", []int64{1, int64(len(prefix))}, prefix)
	if err != nil {
		return nil, err
	}
	encIn, err := floatTensorInput("encoder_hidden_states", enc.Shape, enc.Values)
	if err != nil {
		return nil, err
	}
	inputs := []base.InferInput{idsIn, encIn}
	return c.triton.Infer(ctx, c.cfg.DecoderModel, c.cfg.ModelVersion, inputs, outAllOutputs(), &options.InferOptions{})
}

func (c *Client) decoderWithPast(ctx context.Context, lastToken int64, prev base.InferResult) (base.InferResult, error) {
	idsIn, err := int64TensorInput("input_ids", []int64{1, 1}, []int64{lastToken})
	if err != nil {
		return nil, err
	}
	pasts, err := pastInputsFromResult(prev)
	if err != nil {
		return nil, fmt.Errorf("past kv: %w", err)
	}
	inputs := append([]base.InferInput{idsIn}, pasts...)
	return c.triton.Infer(ctx, c.cfg.DecoderWithPastModel, c.cfg.ModelVersion, inputs, outAllOutputs(), &options.InferOptions{})
}
