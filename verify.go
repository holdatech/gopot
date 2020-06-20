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

	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return err
	}

	var data interface{}

	json.Unmarshal(body, &data)

	headerSignature := r.Header.Get("X-Pot-Signature")

	signedBody, err := CreateSignature(data, secret)
	if err != nil {
		return err
	}

	if headerSignature != signedBody {
		return errors.New("Invalid signature")
	}

	r.Body = ioutil.NopCloser(&buf)

	return nil
}
