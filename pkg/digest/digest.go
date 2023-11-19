package digest

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 will hash the current string using SHA256 and return it as a string
func SHA256(inputString string) string {
	// Create a new SHA-256 hash
	hasher := sha256.New()

	// Write the string to the hash
	hasher.Write([]byte(inputString))

	// Get the finalized hash result as a byte slice
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
