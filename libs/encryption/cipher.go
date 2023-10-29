package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"log"
	"net"
)

var (
	ErrorCipherNotSupported = errors.New("cipher not supported")
	ErrorCipherKeySize      = errors.New("cipher key size error")
)

type Cipher interface {
	KeySize() int
	SaltSize() int
	Cipher() cipher.AEAD
	Secure(conn net.Conn) net.Conn
}

type genericCipher struct {
	key    []byte
	cipher cipher.AEAD
}

func (c *genericCipher) KeySize() int { return len(c.key) }
func (c *genericCipher) SaltSize() int {
	if ks := c.KeySize(); ks > 16 {
		return ks
	}
	return 16
}
func (c *genericCipher) Cipher() cipher.AEAD {
	return c.cipher
}
func (c *genericCipher) Secure(conn net.Conn) net.Conn {
	if c.cipher == nil {
		return conn
	}
	return Secure(conn, c)
}

func checkKeySize(key string, size int) error {
	if len(key) != size {
		return ErrorCipherKeySize
	}
	return nil
}

func CreateCipher(method, key string) (Cipher, error) {
	switch method {
	case "none":
		return &genericCipher{
			key:    []byte(key),
			cipher: nil,
		}, nil
	case "aes-128-gcm":
		err := checkKeySize(key, 16)
		if err != nil {
			log.Printf("cannot create cipher: %s\n", err)
			return nil, err
		}
		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			log.Printf("cannot create cipher: %s\n", err)
			return nil, err
		}
		aesgcm, err := cipher.NewGCM(block)
		return &genericCipher{key: []byte(key), cipher: aesgcm}, nil
	}
	return nil, ErrorCipherNotSupported
}
