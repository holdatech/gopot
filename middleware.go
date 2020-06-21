package gopot

import (
	"net/http"
)

// SignatureVerifierMiddleware creates a middleware with the provided secret to
// verify the signatures in the incoming requests
func SignatureVerifierMiddleware(secret []byte) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := VerifySignatureFromRequest(r, secret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
