package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Character sets for password generation
const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

const (
	defaultPasswordLength = 16
	minPasswordLength     = 4
	maxPasswordLength     = 128
)

// GeneratePassword creates a secure random password with the specified settings
func GeneratePassword(length int, includeLower, includeUpper, includeNumbers, includeSymbols bool) (string, error) {
	if length < minPasswordLength {
		return "", fmt.Errorf("password length must be at least %d", minPasswordLength)
	}
	if length > maxPasswordLength {
		return "", fmt.Errorf("password length must not exceed %d", maxPasswordLength)
	}

	var charset string
	if includeLower {
		charset += lowercase
	}
	if includeUpper {
		charset += uppercase
	}
	if includeNumbers {
		charset += numbers
	}
	if includeSymbols {
		charset += symbols
	}

	if charset == "" {
		return "", fmt.Errorf("at least one character type must be selected")
	}

	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charset[n.Int64()]
	}

	return string(password), nil
}
