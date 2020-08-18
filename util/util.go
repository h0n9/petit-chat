package util

import (
	"crypto/sha256"

	"github.com/h0n9/petit-chat/types"
)

func ToSHA256(data []byte) types.Hash {
	return sha256.Sum256(data)
}
