package signers

import "testing"

func TestNewRSASigner(t *testing.T) {
	signer := NewRsaSigner()
	if signer == nil {
		t.Errorf("NewRsaSigner() returned nil, want non-nil")
	}
}

func TestGenerateKey_RSA_Success(t *testing.T) {
	signer := NewRsaSigner()
	key, err := signer.GenerateKey()
	if err != nil {
		t.Errorf("GenerateKey() returned error: %v, want nil", err)
	}
	if key == nil {
		t.Errorf("GenerateKey() returned nil, want non-nil crypto.Signer")
	}
}
