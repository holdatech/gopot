package gopot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var json jsoniter.API

func init() {
	// Initialize the jsoniter configuration with the custom ascii escape encoder
	jsoniter.RegisterTypeEncoderFunc("string", asciiEncode, asciiIsEmpty)
	config := jsoniter.Config{
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}

	// Freeze the jsoniter API
	json = config.Froze()
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

// CreateSignature creates a pot signature with the given secret
func CreateSignature(d interface{}, secret []byte) (string, error) {
	jdata, err := Marshal(d)

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
