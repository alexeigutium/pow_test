package utils

import (
	"bytes"
	"crypto/sha1"
	"math/bits"
)

// IsValidHash checks if provided hash has required amount of leading zeros
func IsValidHash(zerosLength int, hash []byte) bool {
	for _, b := range hash {
		if zerosLength <= 0 {
			return true
		}
		zeros := bits.LeadingZeros8(uint8(b))

		if zeros == 8 {
			zerosLength -= 8
			continue
		}
		return zeros >= zerosLength
	}

	return false
}

// IsCorrectHash checks if provided hash is correct
func IsCorrectHash(zerosLength int, hash []byte, data []byte) bool {
	if !IsValidHash(zerosLength, hash) {
		return false
	}

	hasher := sha1.New()
	hasher.Write([]byte(data))
	hashed := hasher.Sum(nil)
	return bytes.Equal(hash, hashed)
}
