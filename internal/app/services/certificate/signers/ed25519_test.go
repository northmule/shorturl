package signers

import "testing"

func TestNewEd25519Signer(t *testing.T) {
	signer := NewEd25519Signer()
	if signer == nil {
		t.Errorf("NewEd25519Signer() returned nil, want non-nil")
	}
}

func TestGenerateKey_Ed25519_Success(t *testing.T) {
	signer := NewEd25519Signer()
	key, err := signer.GenerateKey()
	if err != nil {
		t.Errorf("GenerateKey() returned error: %v, want nil", err)
	}
	if key == nil {
		t.Errorf("GenerateKey() returned nil, want non-nil crypto.Signer")
	}
}
