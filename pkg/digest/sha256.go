package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func SHA256(inputString string) string {
	// Create a new SHA-256 hash
	hasher := sha256.New()

	// Write the string to the hash
	hasher.Write([]byte(inputString))

	// Get the finalized hash result as a byte slice
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func FileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	// Copy the file content into the hash object
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Get the hash sum as a byte slice
	hashInBytes := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
