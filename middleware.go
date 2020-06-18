package gopot

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func VerifySignatureFromRequest(r *http.Request, secret []byte) error {
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	body := ioutil.ReadAll(tee)

	if body != secret {
		return errors.New("Invalid signature")
	}

	r.Body = buf

	return nil
}

func SignatureVerifierMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	})
}
