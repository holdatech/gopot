package gopot

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// VerifySignature verifies the signature with the provided public key
func VerifySignature(d interface{}, signature string, key *rsa.PublicKey) error {
	hash, err := calculateHash(d)
	if err != nil {
		return err
	}

	sg, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sg)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

	return nil
}

// VerifySignatureFromRequest can be used to verify the signature in the http request.
// It signs the request body and compares it to the signature provided in the header
func VerifySignatureFromRequest(r *http.Request, key *rsa.PublicKey) error {
	if key == nil {
		return ErrNoSecret
	}
	if r.Body == nil {
		return ErrNoBody
	}

	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return err
	}

	var data interface{}

	json.Unmarshal(body, &data)

	headerSignature := r.Header.Get("X-Pot-Signature")

	if headerSignature == "" {
		return ErrNoSignature
	}

	err = VerifySignature(data, headerSignature, key)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(&buf)

	return nil
}
