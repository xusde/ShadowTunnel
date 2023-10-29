package old

import (
	"testing"
)

func TestEncryptionAndDecryption(t *testing.T) {
	encodedMsg := Encrypt("12345678901234567890123456789012", []byte("hello world"))
	decryptedMsg := Decrypt("12345678901234567890123456789012", encodedMsg)
	if string(decryptedMsg) != "hello world" {
		t.Error("encryption and decryption test failed")
	}
}
