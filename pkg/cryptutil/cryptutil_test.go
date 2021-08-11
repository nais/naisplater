package cryptutil_test

import (
	"bytes"
	"encoding/base64"
	"github.com/nais/naisplater/pkg/cryptutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecryptWithPassword(t *testing.T) {
	plaintext := "plaintext"
	password := "secure"

	ciphertext, err := cryptutil.EncryptWithPassword(plaintext, password)
	assert.NoError(t, err)
	t.Logf("Ciphertext: %s", ciphertext)

	decrypted, err := cryptutil.DecryptWithPassword(ciphertext, password)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptIfPlaintext(t *testing.T) {
	plaintext := "plaintext"
	password := "secure"

	ciphertext, err := cryptutil.EncryptIfPlaintext(plaintext, password)
	assert.NoError(t, err)

	doubleEncrypt, err := cryptutil.EncryptIfPlaintext(ciphertext, password)
	assert.NoError(t, err)
	assert.Equal(t, ciphertext, doubleEncrypt)
}

func TestEncryptIfPlaintextWithWrongKey(t *testing.T) {
	plaintext := "plaintext"
	password := "secure"
	newPassword := "moresecure"

	ciphertext, err := cryptutil.EncryptIfPlaintext(plaintext, password)
	assert.NoError(t, err)

	_, err = cryptutil.EncryptIfPlaintext(ciphertext, newPassword)
	assert.Error(t, err)
}

func TestEncryptIfPlaintextWithGarbage(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	enc.Write(cryptutil.Magic)
	enc.Write([]byte("132146584684516984565146654968475168465846854651356546584"))
	password := "secure"

	_, err := cryptutil.EncryptIfPlaintext(buf.String(), password)
	assert.EqualError(t, err, "cipher: message authentication failed")
}
