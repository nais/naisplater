package cryptutil_test

import (
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
