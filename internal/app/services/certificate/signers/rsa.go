package signers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// RsaSigner Реализация алгоритма rsa
type RsaSigner struct {
}

// NewRsaSigner конструктора
func NewRsaSigner() *RsaSigner {
	return &RsaSigner{}
}

// GenerateKey получение crypto.Signer
func (k *RsaSigner) GenerateKey() (crypto.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	return key, nil
}
