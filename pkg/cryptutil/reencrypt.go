package cryptutil

import (
	"strings"
)

// Decrypt data encrypted with OpenSSL, and re-encrypt with secure parameters.
func ReEncrypt(ciphertext, password string) (string, error) {
	plaintext, err := DecryptOpenSSL(ciphertext, password)
	if err != nil {
		return "", err
	}
	plaintext = strings.TrimSpace(plaintext)
	return EncryptWithPassword(plaintext, password)
}
