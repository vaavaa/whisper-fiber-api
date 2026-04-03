import os

import numpy as np
import triton_python_backend_utils as pb_utils
from transformers import WhisperFeatureExtractor

# Должно совпадать с третьей осью input_features у whisper_*_fp16_encoder (30 с при стандартном hop).
ENCODER_TIME_FRAMES = 3000


class TritonPythonModel:
    def initialize(self, args):
        # Имя HF-модели: должно соответствовать ONNX (medium).
        model_id = os.environ.get(
            "WHISPER_FEATURE_EXTRACTOR_MODEL", "openai/whisper-medium"
        )
        self.feature_extractor = WhisperFeatureExtractor.from_pretrained(model_id)

    def execute(self, requests):
        responses = []
        for request in requests:
            in_tensor = pb_utils.get_input_tensor_by_name(request, "AUDIO_PCM")
            if in_tensor is None:
                err = pb_utils.TritonError("вход AUDIO_PCM отсутствует")
                responses.append(pb_utils.InferenceResponse(error=err))
                continue

            audio = in_tensor.as_numpy().astype(np.float32).flatten()
            try:
                feats = self.feature_extractor(
                    audio,
                    sampling_rate=16000,
                    return_tensors="np",
                )
                mel = np.asarray(feats["input_features"], dtype=np.float32)
            except Exception as exc:  # noqa: BLE001
                err = pb_utils.TritonError(f"feature extractor: {exc}")
                responses.append(pb_utils.InferenceResponse(error=err))
                continue

            if mel.ndim != 3 or mel.shape[1] != 80:
                err = pb_utils.TritonError(
                    f"ожидались размерности mel [1,80,T], получено {mel.shape}"
                )
                responses.append(pb_utils.InferenceResponse(error=err))
                continue

            mel = self._pad_or_truncate_time(mel, ENCODER_TIME_FRAMES)
            out = pb_utils.Tensor("MEL_FEATURES", mel)
            responses.append(pb_utils.InferenceResponse(output_tensors=[out]))
        return responses

    def _pad_or_truncate_time(self, mel: np.ndarray, target_t: int) -> np.ndarray:
        _, _, t = mel.shape
        pad_val = float(getattr(self.feature_extractor, "padding_value", 0.0))
        if t > target_t:
            return mel[:, :, :target_t].astype(np.float32)
        if t < target_t:
            pad_w = target_t - t
            return np.concatenate(
                [
                    mel.astype(np.float32),
                    np.full((1, 80, pad_w), pad_val, dtype=np.float32),
                ],
                axis=2,
            )
        return mel.astype(np.float32)

    def finalize(self):
        pass
