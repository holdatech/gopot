package gopot

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
)

// VerifySignature verifies the signature with the provided public key
func VerifySignature(d interface{}, signature string, keys []*rsa.PublicKey) error {
	hash, err := calculateHash(d)
	if err != nil {
		return err
	}

	sg, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	valid := false
	for _, key := range keys {
		err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sg)
		if err == nil {
			valid = true
			break
		}
	}

	if !valid {
		return ErrInvalidSignature
	}

	return nil
}

// VerifySignatureFromRequest can be used to verify the signature in the http request.
// It signs the request body and compares it to the signature provided in the header
func VerifySignatureFromRequest(r *http.Request, keys []*rsa.PublicKey) error {
	if len(keys) == 0 {
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

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	headerSignature := r.Header.Get("X-Pot-Signature")

	if headerSignature == "" {
		return ErrNoSignature
	}

	err = VerifySignature(data, headerSignature, keys)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(&buf)

	return nil
}
