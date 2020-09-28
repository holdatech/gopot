package gopot

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
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

// calculateHash returns a calculated hash from the provided data
func calculateHash(d interface{}) (hash [32]byte, err error) {
	jdata, err := Marshal(d)
	if err != nil {
		return
	}

	hash = sha256.Sum256(jdata)
	return
}

// CreateSignature creates a pot signature with the given secret
func CreateSignature(d interface{}, key *rsa.PrivateKey) (string, error) {
	if key == nil {
		return "", ErrNoSecret
	}

	hash, err := calculateHash(d)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), err
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
