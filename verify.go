package gopot

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// VerifySignatureFromRequest can be used to verify the signature in the http request.
// It signs the request body and compares it to the signature provided in the header
func VerifySignatureFromRequest(r *http.Request, secret []byte) error {
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

	signedBody, err := CreateSignature(data, secret)
	if err != nil {
		return err
	}

	if headerSignature != signedBody {
		return ErrInvalidSignature
	}

	r.Body = ioutil.NopCloser(&buf)

	return nil
}
