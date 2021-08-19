package util

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/h0n9/petit-chat/types"
)

func ToSHA256(data []byte) types.Hash {
	return sha256.Sum256(data)
}

func ToBool(boolStr string) (bool, error) {
	switch boolStr {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("only boolean type: 'true' or 'false'")
	}
}

func HasField(name string, s interface{}) bool {
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}
