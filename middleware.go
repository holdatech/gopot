package gopot

import (
	"crypto/rsa"
	"net/http"
)

// SignatureVerifierMiddleware creates a middleware with the provided secret to
// verify the signatures in the incoming requests
func SignatureVerifierMiddleware(pubKey *rsa.PublicKey) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := VerifySignatureFromRequest(r, pubKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
