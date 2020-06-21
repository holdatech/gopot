package gopot

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

var requestTestBody = []byte(`
{
	"parameters":{
		"businessId":"2980005-2",
		"test":"Hei vaan äöääöö \":{}\""
	},
	"productCode":"prh-business-identity-data-product",
	"timestamp":"2020-05-19T14:31:24Z"
}
`)

var errTestError = errors.New("test error")

func TestVerifySignatureFromRequest(t *testing.T) {
	var secret = []byte("P8qNkpXkfLe_OQa_2ydHRgzFR2_GuIoyUoMtf8zcLZ0")

	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	err := VerifySignatureFromRequest(req, secret)
	if err != nil {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(req.Body)

	t.Logf("%s", body)
}

func TestVerifySignatureNoBody(t *testing.T) {
	var secret = []byte("P8qNkpXkfLe_OQa_2ydHRgzFR2_GuIoyUoMtf8zcLZ0")

	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")

	err := VerifySignatureFromRequest(req, secret)
	if err != nil && !errors.Is(err, ErrNoBody) {
		t.Error(err)
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errTestError
}

func TestVerifySignatureErrBody(t *testing.T) {
	var secret = []byte("P8qNkpXkfLe_OQa_2ydHRgzFR2_GuIoyUoMtf8zcLZ0")

	req, _ := http.NewRequest("POST", "/fetch", errReader(0))
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")

	err := VerifySignatureFromRequest(req, secret)
	if err != nil && !errors.Is(err, errTestError) {
		t.Error(err)
	}
}

func TestVerifySignatureNoSecret(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgyvQHyKCXWM=")

	err := VerifySignatureFromRequest(req, nil)
	if err != nil && !errors.Is(err, ErrNoSecret) {
		t.Error(err)
	}
}
