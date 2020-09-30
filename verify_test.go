package gopot

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errTestError
}

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

var secretKey, _ = rsa.GenerateKey(rand.Reader, 4096)
var errTestError = errors.New("test error")

func testSignature() string {
	var data interface{}

	json.Unmarshal(requestTestBody, &data)
	signature, _ := CreateSignature(data, secretKey)

	return signature
}

func TestVerifySignature(t *testing.T) {
	var data interface{}
	json.Unmarshal(requestTestBody, &data)
	signature, _ := CreateSignature(data, secretKey)

	err := VerifySignature(data, signature, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil {
		t.Error(err)
	}
}

func TestVerifySignatureMultipeKeys(t *testing.T) {
	var data interface{}
	json.Unmarshal(requestTestBody, &data)
	signature, _ := CreateSignature(data, secretKey)

	secretKey2, _ := rsa.GenerateKey(rand.Reader, 4096)

	err := VerifySignature(data, signature, []*rsa.PublicKey{&secretKey2.PublicKey, &secretKey.PublicKey})
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	var data interface{}
	json.Unmarshal(requestTestBody, &data)
	signature, _ := CreateSignature(data, secretKey)

	for i := 0; i < b.N; i++ {
		VerifySignature(data, signature, []*rsa.PublicKey{&secretKey.PublicKey})
	}
}

func TestVerifySignatureFromRequest(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", testSignature())
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	err := VerifySignatureFromRequest(req, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(req.Body)

	t.Logf("%s", body)
}

func TestVerifySignatureInvalidSignature(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	err := VerifySignatureFromRequest(req, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil && !errors.Is(err, ErrInvalidSignature) {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(req.Body)

	t.Logf("%s", body)
}

func TestVerifySignatureNoBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")

	err := VerifySignatureFromRequest(req, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil && !errors.Is(err, ErrNoBody) {
		t.Error(err)
	}
}

func TestVerifySignatureErrBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", errReader(0))
	req.Header.Set("X-Pot-Signature", testSignature())

	err := VerifySignatureFromRequest(req, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil && !errors.Is(err, errTestError) {
		t.Error(err)
	}
}

func TestVerifySignatureNoSecret(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))
	req.Header.Set("X-Pot-Signature", testSignature())

	err := VerifySignatureFromRequest(req, nil)
	if err != nil && !errors.Is(err, ErrNoSecret) {
		t.Error(err)
	}
}

func TestVerifySignatureNoSignature(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Body = ioutil.NopCloser(bytes.NewReader(requestTestBody))

	err := VerifySignatureFromRequest(req, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil && !errors.Is(err, ErrNoSignature) {
		t.Error(err)
	}
}
