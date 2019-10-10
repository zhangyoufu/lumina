package lumina

import "errors"

// These errors are used by serialization/deserialization.
var (
	errInvalidPayloadLength = errors.New("invalid payload length")
	errTrailingData         = errors.New("trailing data at end of unmarshal")
	errUnexpectedEncoding   = errors.New("unexpected encoding")
	errTooLong              = errors.New("too long")
	errInvalidValue         = errors.New("invalid value")
)
