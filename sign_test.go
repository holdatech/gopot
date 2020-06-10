package gopot

import (
	"bytes"
	"testing"
	"time"
)

func TestAddSpaces(t *testing.T) {
	cases := []struct {
		in  []byte
		out []byte
	}{
		{
			in:  []byte(`{"hello":"world","foo":{"bar":"baz","testData":"\"äö£€\":{}\""}}`),
			out: []byte(`{"hello": "world","foo": {"bar": "baz","testData": "\"äö£€\":{}\""}}`),
		},
	}

	for _, c := range cases {
		d := JSONAddSpaces(c.in)

		if !bytes.Equal(d, c.out) {
			t.Errorf("Failed to add spaces: in: %s, out: %s, expected: %s", c.in, d, c.out)
		}
	}
}

func TestAsciiString(t *testing.T) {
	cases := []struct {
		in  string
		out []byte
	}{
		{
			in:  "hello",
			out: []byte(`"hello"`),
		},
		{
			in:  "hello & foo € bar",
			out: []byte(`"hello & foo \u20ac bar"`),
		},
		{
			in:  "äöuyå",
			out: []byte(`"\u00e4\u00f6uy\u00e5"`),
		},
		{
			in:  `! "$\#%&'()*,+-./:;<=>?[]{}~|_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`,
			out: []byte(`"! \"$\\#%&'()*,+-./:;<=>?[]{}~|_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"`),
		},
	}

	for _, c := range cases {
		d, _ := json.Marshal(c.in)
		if !bytes.Equal(d, c.out) {
			t.Errorf("Encode ascii string to json: in: %s, out: %s, expected: %s", c.in, d, c.out)
		}
	}
}

func TestSignature(t *testing.T) {
	type parameters struct {
		BusinessID string `json:"businessId"`
		Test       string `json:"test"`
	}

	type testData struct {
		Parameters  parameters `json:"parameters"`
		ProductCode string     `json:"productCode"`
		Timestamp   time.Time  `json:"timestamp"`
	}

	var secret = "P8qNkpXkfLe_OQa_2ydHRgzFR2_GuIoyUoMtf8zcLZ0"

	ts, _ := time.Parse(time.RFC3339, "2020-05-19T14:31:24Z")
	test := &testData{
		Timestamp:   ts,
		ProductCode: "prh-business-identity-data-product",
		Parameters: parameters{
			BusinessID: "2980005-2",
			Test:       "Hei vaan äöääöö \":{}\"",
		},
	}

	signature, _ := CreateSignature(test, []byte(secret))

	if signature != "5t1XQofwg2Uc6j7LnhNz0gvFL0AgJj0sGyvQHyKCXWM=" {
		t.Error("Signature doesn't match")
	}

	t.Log(signature)
}
