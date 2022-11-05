package constant

import "errors"

// Configuration error
var (
	ErrMissingToken = errors.New("[Configuration] Missing GitHub access token")
)

// Metadata error
var (
	ErrMalformedMetadata = errors.New("[Metadata] Malformed repository metadata")
)

// Event error
var (
	ErrEventFileRead  = errors.New("[Event] Failed to read event file")
	ErrEventFileParse = errors.New("[Event] Failed to parse event file")
)
