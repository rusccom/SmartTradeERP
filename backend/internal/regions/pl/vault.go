package pl

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

func binding(tenant string, s Seller) []byte {
	return []byte("ksef:" + tenant + ":" + s.Environment + ":" + s.NIP)
}

func (h *Handler) seal(tenant string, u profileUpdate) ([]byte, error) {
	if u.Token == "" {
		return nil, nil
	}
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(h.key))
	if err != nil || len(key) != chacha20poly1305.KeySize {
		return nil, problem("KSeF credential storage requires a server encryption key.")
	}
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, cipher.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	return cipher.Seal(nonce, nonce, []byte(u.Token), binding(tenant, u.Seller)), nil
}

func (h *Handler) decrypt(tenant string, s Seller, stored []byte) (string, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(h.key))
	if err != nil || len(key) != chacha20poly1305.KeySize {
		return "", problem("KSeF credential storage is unavailable.")
	}
	box, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err
	}
	if len(stored) < box.NonceSize()+box.Overhead() {
		return "", problem("Configure the seller's KSeF token first.")
	}
	plain, err := box.Open(nil, stored[:box.NonceSize()], stored[box.NonceSize():], binding(tenant, s))
	if err != nil {
		return "", problem("KSeF credentials cannot be opened. Re-enter the token in seller settings.")
	}
	return string(plain), nil
}
