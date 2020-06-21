package gopot

import "errors"

var ErrNoBody = errors.New("No body")
var ErrInvalidSignature = errors.New("Invalid signature")
var ErrNoSecret = errors.New("No secret provided")
