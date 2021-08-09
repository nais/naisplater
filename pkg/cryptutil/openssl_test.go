// +build integration
//
// requires openssl installed on machine

package cryptutil_test

import (
	"github.com/nais/naisplater/pkg/cryptutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecryptOpenSSL(t *testing.T) {
	plaintext := "plaintext"
	ciphertext := "U2FsdGVkX19JGpnp/eKAP1tK0hUr+8jHtPeiw7rKFng="
	key := "secure"

	decrypted, err := cryptutil.DecryptOpenSSL(ciphertext, key)

	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}