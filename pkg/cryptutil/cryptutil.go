package cryptutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
	"strings"
)

const iterations = 10000
const saltlen = 16
const keylen = 32

var Magic = []byte("CRYPT")
var ErrNotEncrypted = errors.New("not an encrypted value")

func pbkdf(key, salt []byte) []byte {
	return pbkdf2.Key(key, salt, iterations, keylen, sha256.New)
}

// Encrypt and base64-encode data using aes-256-gcm and a key derived from a password hashed with PBKDF2.
// Output is a base64-encoded string with 16 bytes of PBKDF2 salt, 12 bytes of iv, followed by ciphertext.
func EncryptWithPassword(plaintext string, password string) (string, error) {
	salt, err := randomBytes(saltlen)
	if err != nil {
		return "", err
	}
	key := pbkdf([]byte(password), salt)

	ciphertext, err := Encrypt([]byte(plaintext), key)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	defer enc.Close()

	_, err = enc.Write(Magic)
	if err != nil {
		return "", err
	}
	_, err = enc.Write(salt)
	if err != nil {
		return "", err
	}
	_, err = enc.Write(ciphertext)
	if err != nil {
		return "", err
	}

	enc.Close()

	return buf.String(), nil
}

// Decrypts a base64-encoded ciphertext encrypted with EncryptWithPassword.
// Returns the plaintext as string.
func DecryptWithPassword(ciphertext string, password string) (string, error) {
	r := strings.NewReader(ciphertext)
	dec := base64.NewDecoder(base64.StdEncoding, r)

	magic, err := readExactly(dec, len(Magic))
	if err != nil {
		return "", err
	}
	if bytes.Compare(magic, Magic) != 0 {
		return "", ErrNotEncrypted
	}

	salt, err := readExactly(dec, saltlen)
	if err != nil {
		return "", err
	}

	key := pbkdf([]byte(password), salt)

	encrypted, err := ioutil.ReadAll(dec)
	if err != nil {
		return "", err
	}

	plaintext, err := Decrypt(encrypted, key)

	return string(plaintext), err
}

func readExactly(r io.Reader, length int) ([]byte, error) {
	data := make([]byte, length)
	nread, err := io.ReadAtLeast(r, data, length)
	if err != nil {
		return nil, err
	}
	if nread != length {
		return nil, fmt.Errorf("too much data read")
	}
	return data, nil
}
