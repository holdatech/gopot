package gopot

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestVerifySignatureFromRequest(t *testing.T) {
	var secret = []byte("P8qNkpXkfLe_OQa_2ydHRgzFR2_GuIoyUoMtf8zcLZ0")

	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("X-Pot-Signature", "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"parameters":{"businessId":"2980005-2","test":"Hei vaan äöääöö \":{}\""},"productCode":"prh-business-identity-data-product","timestamp":"2020-05-19T14:31:24Z"}`)))

	err := VerifySignatureFromRequest(req, secret)
	if err != nil {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(req.Body)

	t.Logf("%s", body)
}
