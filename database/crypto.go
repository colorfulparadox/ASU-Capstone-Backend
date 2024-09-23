package database

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func GenerateUUID() string {
	dateIssued := time.Now().Unix()
	expires := time.Now().AddDate(0, 0, 7).Unix()
	num, err := rand.Int(rand.Reader, big.NewInt(int64(10000)))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return ""
	}

	token := uuid.New().String() + "/" + strconv.FormatInt(expires+dateIssued^(num.Int64()), 16) + "/" + GenerateRandomStringURLSafe(16)
	return token

}

func Randomize_auth_token(auth_token string) {
	randomize_auth_token(auth_token)
}

// Hashes the users password for storage
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return string(hashedPassword)
	}

	fmt.Printf("Hashed Password: %s\n", string(hashedPassword))

	return string(hashedPassword)
}

// Verifies the password against the stored hash
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			fmt.Println("Invalid password")
		} else {
			fmt.Println("Error verifying password:", err)
		}
		return false
	}
	return true
}
