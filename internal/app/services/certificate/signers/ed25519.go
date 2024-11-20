package signers

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
)

// Ed25519Signer Реализация алгоритма Ed25519
type Ed25519Signer struct {
}

// NewEd25519Signer конструктора
func NewEd25519Signer() *Ed25519Signer {
	return &Ed25519Signer{}
}

// GenerateKey получение crypto.Signer
func (k *Ed25519Signer) GenerateKey() (crypto.Signer, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return key, nil
}
