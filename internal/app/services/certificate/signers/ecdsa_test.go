package signers

import (
	"crypto"
	"errors"
	"testing"
)

// BadSigner реализация с ошибкой
type BadSigner struct {
}

// GenerateKey метод возвращает ошибку
func (b *BadSigner) GenerateKey() (crypto.Signer, error) {
	return nil, errors.New("GenerateKey error")
}

func TestNewEcdsaSigner(t *testing.T) {
	signer := NewEcdsaSigner()
	if signer == nil {
		t.Errorf("NewEcdsaSigner() returned nil, want non-nil")
	}
}

func TestGenerateKey_Ecdsa_Success(t *testing.T) {
	signer := NewEcdsaSigner()
	key, err := signer.GenerateKey()
	if err != nil {
		t.Errorf("GenerateKey() returned error: %v, want nil", err)
	}
	if key == nil {
		t.Errorf("GenerateKey() returned nil, want non-nil crypto.Signer")
	}
}

func TestGenerateKey_Failure(t *testing.T) {
	signer := new(BadSigner)
	_, err := signer.GenerateKey()
	if err == nil {
		t.Errorf("GenerateKey() returned nil error, want non-nil")
	}
}
