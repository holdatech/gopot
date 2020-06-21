package gopot

import "errors"

var (
	// ErrNoBody is returned when there's no body provided to the signature function
	ErrNoBody = errors.New("No body")
	// ErrInvalidSignature is returned when the signatures don't match
	ErrInvalidSignature = errors.New("Invalid signature")
	// ErrNoSecret is returned when there's no secret provided
	ErrNoSecret = errors.New("No secret provided")
	// ErrNoSignature is returned when there's no signature found in the request headers
	ErrNoSignature = errors.New("No signature provided")
)
