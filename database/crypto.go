package database

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generaterandombytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generaterandomstringURLsafe(n int) (string, error) {
	b, err := generaterandombytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomStringURLSafe(n int) string {
	token, err := generaterandomstringURLsafe(n)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}
	return token
}
