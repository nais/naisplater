package cryptutil

import (
	"os/exec"
	"strings"
)

func DecryptOpenSSL(ciphertext, password string) (string, error) {
	cmd := exec.Command("openssl", "enc", "-d", "-aes-256-cbc", "-a", "-A", "-k", password)
	cmd.Stdin = strings.NewReader(ciphertext)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
