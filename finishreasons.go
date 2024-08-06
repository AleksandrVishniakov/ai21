package ai21

type FinishReason string

const (
	// FinishReasonStop indicates that the response ended naturally as a complete answer
	// (due to end-of-sequence token) or because the model generated a stop sequence
	// provided in the request.
	FinishReasonStop FinishReason = "stop"

	// FinishReasonLength indicates that the response ended by reaching max_tokens.
	FinishReasonLength FinishReason = "length"
)
