package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

var debugEnabled bool

func HashObject(obj interface{}) (string, error) {
	// Serialize the object to JSON
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	// Hash the serialized data
	hasher := sha256.New()
	hasher.Write(b)
	hash := hasher.Sum(nil)

	// Return the hexadecimal encoding of the hash
	return fmt.Sprintf("%x", hash), nil
}

func Debugf(str string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf(str, args...)
	}
}

func Debugln(str ...interface{}) {
	if debugEnabled {
		fmt.Println(str...)
	}
}

func SetDebugEnabled(enabled bool) {
	debugEnabled = enabled
}
