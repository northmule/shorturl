package signers

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// EcdsaSigner Реализация алгоритма ecdsa
type EcdsaSigner struct {
}

// NewEcdsaSigner конструктора
func NewEcdsaSigner() *EcdsaSigner {
	return &EcdsaSigner{}
}

// GenerateKey получение crypto.Signer
func (k *EcdsaSigner) GenerateKey() (crypto.Signer, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return key, nil
}
