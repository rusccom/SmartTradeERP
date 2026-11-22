package basecourier

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

type vault struct{ cipher.AEAD }

func newVault(encoded string) vault {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return vault{}
	}
	box, err := chacha20poly1305.NewX(key)
	if err != nil {
		return vault{}
	}
	return vault{box}
}

func keyBinding(tenantID, email string) []byte {
	return []byte("base-courier:" + tenantID + ":" + email)
}

func (v vault) seal(tenantID, email, key string) ([]byte, error) {
	if v.AEAD == nil {
		return nil, errVault
	}
	nonce := make([]byte, v.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, errVault
	}
	return v.Seal(nonce, nonce, []byte(key), keyBinding(tenantID, email)), nil
}

func (v vault) open(tenantID, email string, sealed []byte) (string, error) {
	if v.AEAD == nil || len(sealed) < v.NonceSize()+v.Overhead() {
		return "", errVault
	}
	n := v.NonceSize()
	plain, err := v.Open(nil, sealed[:n], sealed[n:], keyBinding(tenantID, email))
	if err != nil {
		return "", errVault
	}
	return string(plain), nil
}
