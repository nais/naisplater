package cryptutil

// Decrypt data encrypted with OpenSSL, and re-encrypt with secure parameters.
func ReEncrypt(ciphertext, password string) (string, error) {
	plaintext, err := DecryptOpenSSL(ciphertext, password)
	if err != nil {
		return "", err
	}

	return EncryptWithPassword(plaintext, password)
}
