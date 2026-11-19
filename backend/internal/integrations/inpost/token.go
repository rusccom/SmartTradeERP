package inpost

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"strings"
)

type credential struct {
	plain     string
	encrypted []byte
}

func newTokenCipher(encoded string) cipher.AEAD {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil || len(key) != 32 {
		return nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	return aead
}

func sealToken(aead cipher.AEAD, token string, binding []byte) (credential, error) {
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return credential{}, errVault
	}
	return credential{token, aead.Seal(nonce, nonce, []byte(token), binding)}, nil
}

func openToken(aead cipher.AEAD, encrypted, binding []byte) (credential, error) {
	size := aead.NonceSize()
	if len(encrypted) < size+aead.Overhead() {
		return credential{}, errVault
	}
	plain, err := aead.Open(nil, encrypted[:size], encrypted[size:], binding)
	if err != nil {
		return credential{}, errVault
	}
	return credential{string(plain), encrypted}, nil
}
