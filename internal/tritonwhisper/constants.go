package tritonwhisper

// Совпадает с openai/whisper-medium config.json (vocab, forced_decoder_ids).
const (
	DecoderLayers = 24
	EOSTokenID           = 50257
	DecoderStartTokenID  = 50258
	MaxDecoderPositions  = 448 // max_target_positions в config.json whisper-medium
	DefaultAudioSampleHz = 16000
)

// DefaultDecoderPrefix — forced_decoder_ids из config.json: после decoder_start идут токены задачи (transcribe).
// Для другого языка замените на соответствующие ID из generation_config.json / tokenizer.
var DefaultDecoderPrefix = []int64{DecoderStartTokenID, 50259, 50359, 50363}
