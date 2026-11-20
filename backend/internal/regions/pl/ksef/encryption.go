package ksef

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"slices"
	"time"
)

type certificate struct {
	Data  string    `json:"certificate"`
	ID    string    `json:"publicKeyId"`
	From  time.Time `json:"validFrom"`
	Until time.Time `json:"validTo"`
	Usage []string  `json:"usage"`
}

func selectCert(certs []certificate, usage string) (certificate, error) {
	now := time.Now()
	for _, c := range certs {
		if slices.Contains(c.Usage, usage) && !now.Before(c.From) && now.Before(c.Until) {
			return c, nil
		}
	}
	return certificate{}, errors.New("KSeF has no currently valid encryption certificate")
}

func (c certificate) wrap(plain []byte) (string, error) {
	der, err := base64.StdEncoding.DecodeString(c.Data)
	if err != nil {
		return "", err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return "", err
	}
	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("unexpected KSeF certificate type")
	}
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, plain, nil)
	return base64.StdEncoding.EncodeToString(enc), err
}

type material struct{ key, iv []byte }

func newMaterial() (material, error) {
	m := material{make([]byte, 32), make([]byte, aes.BlockSize)}
	if _, err := rand.Read(m.key); err != nil {
		return m, err
	}
	_, err := rand.Read(m.iv)
	return m, err
}

func digest(b []byte) string { v := sha256.Sum256(b); return base64.StdEncoding.EncodeToString(v[:]) }

func (m material) invoice(xml []byte) (map[string]any, error) {
	block, err := aes.NewCipher(m.key)
	if err != nil {
		return nil, err
	}
	pad := aes.BlockSize - len(xml)%aes.BlockSize
	plain := append(bytes.Clone(xml), bytes.Repeat([]byte{byte(pad)}, pad)...)
	cipher.NewCBCEncrypter(block, m.iv).CryptBlocks(plain, plain)
	return map[string]any{"invoiceHash": digest(xml), "invoiceSize": len(xml), "encryptedInvoiceHash": digest(plain),
		"encryptedInvoiceSize": len(plain), "encryptedInvoiceContent": base64.StdEncoding.EncodeToString(plain)}, nil
}
