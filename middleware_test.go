package gopot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestSignatureVerifierMiddleware(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SignatureVerifierMiddleware(&secretKey.PublicKey))
	router.Post("/fetch", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	})

	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", testSignature())
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status: Wanted: %d, Got: %d", http.StatusOK, status)
	}

	requestBody, _ := ioutil.ReadAll(rr.Body)

	if !bytes.Equal(requestBody, requestTestBody) {
		t.Error("invalid request body")
	}
}

func TestWrongStatusSignatureVerifierMiddleware(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SignatureVerifierMiddleware(&secretKey.PublicKey))
	router.Post("/fetch", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	})

	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCX=")
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong status")
	}
}
