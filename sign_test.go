package gopot

import (
	"bytes"
	"crypto/rsa"
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

type marshalTest struct {
	Foo   string       `json:"foo"`
	Hello string       `json:"hello,omitempty"`
	Nil   *string      `json:"nil"`
	World string       `json:"world,omitempty"`
	X     float32      `json:"x"`
	Y     float32      `json:"y"`
	Data  *marshalTest `json:"data,omitempty"`
}

var benchmarkMarshalTestData = marshalTest{
	Foo:   "bar",
	Hello: "baz",
	World: "hello",
	X:     1234123.12341234,
	Y:     12323.12341234,
}

func BenchmarkMarshal(b *testing.B) {
	data := &benchmarkMarshalTestData
	dd1 := *data
	dd2 := *data

	data.Data = &dd1
	data.Data.Data = &dd2

	for i := 0; i < b.N; i++ {
		Marshal(data)
	}
}

func BenchmarkMarshalWithoutSpaces(b *testing.B) {
	data := &benchmarkMarshalTestData

	dd1 := *data
	dd2 := *data

	data.Data = &dd1
	data.Data.Data = &dd2

	for i := 0; i < b.N; i++ {
		MarshalWithoutSpaces(data)
	}
}

func BenchmarkCalculateHash(b *testing.B) {
	type testData struct {
		Foo   string  `json:"foo"`
		Hello string  `json:"hello,omitempty"`
		Nil   *string `json:"nil"`
		World string  `json:"world,omitempty"`
	}

	data := &testData{
		Foo:   "bar",
		Hello: "baz",
		World: "hello",
	}

	for i := 0; i < b.N; i++ {
		calculateHash(data)
	}
}

func BenchmarkCreateSignature(b *testing.B) {
	type testData struct {
		Foo   string  `json:"foo"`
		Hello string  `json:"hello,omitempty"`
		Nil   *string `json:"nil"`
		World string  `json:"world,omitempty"`
	}

	data := &testData{
		Foo:   "bar",
		Hello: "baz",
		World: "hello",
	}

	for i := 0; i < b.N; i++ {
		CreateSignature(data, secretKey)
	}
}

func TestAsciiString(t *testing.T) {
	cases := []struct {
		in  string
		out []byte
	}{
		{
			in:  "",
			out: []byte(`""`),
		},
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
			t.Errorf("Encode ascii string to json: in: `%s`, out: `%s`, expected: `%s`", c.in, d, c.out)
		}
	}
}

func TestNilString(t *testing.T) {
	type testData struct {
		Foo   string  `json:"foo"`
		Hello string  `json:"hello,omitempty"`
		Nil   *string `json:"nil"`
		World string  `json:"world,omitempty"`
	}

	d, _ := Marshal(&testData{
		Foo:   "bar",
		Hello: "",
		World: "world",
	})

	expected := []byte(`{"foo": "bar","nil": null,"world": "world"}`)

	if !bytes.Equal(d, expected) {
		t.Errorf("Encode ascii string to json: out: `%s`, expected: `%s`", d, expected)
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

	ts, _ := time.Parse(time.RFC3339, "2020-05-19T14:31:24Z")
	test := &testData{
		Timestamp:   ts,
		ProductCode: "prh-business-identity-data-product",
		Parameters: parameters{
			BusinessID: "2980005-2",
			Test:       "Hei vaan äöääöö \":{}\"",
		},
	}

	signature, _ := CreateSignature(test, secretKey)

	err := VerifySignature(test, signature, []*rsa.PublicKey{&secretKey.PublicKey})
	if err != nil {
		t.Errorf("Signature doesn't match: %e", err)
	}

	t.Log(signature)
}
