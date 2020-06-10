package gopot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strconv"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var json jsoniter.API

func init() {
	jsoniter.RegisterTypeEncoderFunc("string", asciiEncode, asciiIsEmpty)
	conf := jsoniter.Config{
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}
	json = conf.Froze()
}

func asciiEncode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	str := *(*string)(ptr)
	stream.WriteRaw(strconv.QuoteToASCII(str))
}

func asciiIsEmpty(ptr unsafe.Pointer) bool {
	if *(*string)(ptr) == "" {
		return true
	}
	return false
}

// AsciiString can be used to escape the utf-8 charracters in the json output
type AsciiString string

func (s AsciiString) MarshalJSON() ([]byte, error) {
	res := []byte(strconv.QuoteToASCII(string(s)))
	return res, nil
}

// CreateSignature creates a pot signature with the given secret
func CreateSignature(d interface{}, secret []byte) (string, error) {
	jdata, err := Marshal(d)

	log.Printf("payload to be signed: %s", jdata)

	// Sign the payload
	hash := hmac.New(sha256.New, secret)
	hash.Write(jdata)

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), err
}

// Marshal marshals json with the POT separators added.
func Marshal(d interface{}) ([]byte, error) {
	jdata, err := json.Marshal(d)
	jdata = JSONAddSpaces(jdata)

	return jdata, err
}

// JSONAddSpaces ands spaces after the value declarations in json
func JSONAddSpaces(src []byte) []byte {
	var res []byte
	isEscaped := false
	isValue := false
	for _, b := range src {
		res = append(res, b)
		if !isEscaped && b == '"' {
			isValue = !isValue
		}

		if b == ':' && !isValue {
			res = append(res, ' ')
		}
		if b == '\\' && !isEscaped {
			isEscaped = true
			continue
		}
		if isEscaped {
			isEscaped = false
		}
	}

	return res
}
