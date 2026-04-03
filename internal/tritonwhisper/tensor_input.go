package tritonwhisper

import (
	"fmt"

	"github.com/Trendyol/go-triton-client/base"
	tritongrpc "github.com/Trendyol/go-triton-client/client/grpc"
)

func floatTensorInput(name string, shape []int64, data []float32) (base.InferInput, error) {
	in := tritongrpc.NewInferInput(name, "FP32", shape, nil)
	if err := in.SetData(data, true); err != nil {
		return nil, err
	}
	return in, nil
}

func int64TensorInput(name string, shape []int64, data []int64) (base.InferInput, error) {
	in := tritongrpc.NewInferInput(name, "INT64", shape, nil)
	if err := in.SetData(data, true); err != nil {
		return nil, err
	}
	return in, nil
}

func pastInputsFromResult(prev base.InferResult) ([]base.InferInput, error) {
	inputs := make([]base.InferInput, 0, DecoderLayers*4)
	for i := range DecoderLayers {
		for _, suffix := range []string{".decoder.key", ".decoder.value", ".encoder.key", ".encoder.value"} {
			outName := fmt.Sprintf("present.%d%s", i, suffix)
			inName := fmt.Sprintf("past_key_values.%d%s", i, suffix)
			shape, err := prev.GetShape(outName)
			if err != nil {
				return nil, fmt.Errorf("shape %s: %w", outName, err)
			}
			data, err := prev.AsFloat32Slice(outName)
			if err != nil {
				return nil, fmt.Errorf("values %s: %w", outName, err)
			}
			inp, err := floatTensorInput(inName, shape, data)
			if err != nil {
				return nil, fmt.Errorf("input %s: %w", inName, err)
			}
			inputs = append(inputs, inp)
		}
	}
	return inputs, nil
}
