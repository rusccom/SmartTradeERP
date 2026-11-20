package connection

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type Vault struct{ cipher.AEAD }

func NewVault(encoded string) Vault {
	key, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(key) != 32 {
		return Vault{}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return Vault{}
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return Vault{}
	}
	return Vault{aead}
}

func (v Vault) Encrypt(binding string, value any) ([]byte, error) {
	if v.AEAD == nil {
		return nil, Setup
	}
	plain, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, v.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	return v.Seal(nonce, nonce, plain, []byte(binding)), nil
}

func (v Vault) Decrypt(binding string, data []byte, dest any) error {
	if v.AEAD == nil {
		return Setup
	}
	if len(data) < v.NonceSize()+v.Overhead() {
		return Stored
	}
	plain, err := v.Open(nil, data[:v.NonceSize()], data[v.NonceSize():], []byte(binding))
	if err != nil || json.Unmarshal(plain, dest) != nil {
		return Stored
	}
	return nil
}
