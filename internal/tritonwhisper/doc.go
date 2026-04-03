// Package tritonwhisper — скелет gRPC-клиента для Triton: whisper_ensemble (PCM→энкодер)
// и жадный цикл whisper_medium_fp16_decoder / _with_past.
//
// Текст из token ID здесь не собирается: подключите tokenizer.json (например github.com/daulet/tokenizers)
// или декодируйте ID во внешнем сервисе. Префикс декодера (DefaultDecoderPrefix) взят из config.json
// openai/whisper-medium; для других языков замените forced_decoder_ids.
package tritonwhisper
