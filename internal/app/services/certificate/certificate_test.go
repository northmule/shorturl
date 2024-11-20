package certificate

import (
	"crypto"
	"errors"
	"os"
	"testing"

	"github.com/northmule/shorturl/internal/app/services/certificate/signers"
)

// BadSigner реализация с ошибкой
type BadSigner struct {
}

// GenerateKey метод возвращает ошибку
func (b *BadSigner) GenerateKey() (crypto.Signer, error) {
	return nil, errors.New("GenerateKey error")
}

func TestInitSelfSigned(t *testing.T) {
	tests := []struct {
		name   string
		signer KeyGenerator
	}{
		{
			name:   "Rsa",
			signer: signers.NewRsaSigner(),
		},
		{
			name:   "Ecdsa",
			signer: signers.NewEcdsaSigner(),
		},
		{
			name:   "Ed25519",
			signer: signers.NewEd25519Signer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			c := NewCertificate(tt.signer)

			err = c.InitSelfSigned()
			if err != nil {
				t.Fatalf("InitSelfSigned() error = %v", err)
			}

			_, err = os.Stat(c.CertPath())
			if os.IsNotExist(err) {
				t.Errorf("Cert file does not exist")
			}

			_, err = os.Stat(c.KeyPath())
			if os.IsNotExist(err) {
				t.Errorf("Key file does not exist")
			}

			err = os.Remove(c.CertPath())
			if err != nil {
				return
			}
			err = os.Remove(c.KeyPath())
			if err != nil {
				t.Fatal("os.Remove")
			}
		})
	}
}

func TestInitSelfSigned_BadSigner(t *testing.T) {
	var err error
	rsaSigner := new(BadSigner)
	c := NewCertificate(rsaSigner)

	err = c.InitSelfSigned()
	if err == nil {
		t.Fatalf("expected InitSelfSigned() error = %v", err)
	}

}

func BenchmarkCertificate_InitSelfSigned(b *testing.B) {
	tests := []struct {
		name   string
		signer KeyGenerator
	}{
		{
			name:   "Benchmark_Rsa",
			signer: signers.NewRsaSigner(),
		},
		{
			name:   "Benchmark_Ecdsa",
			signer: signers.NewEcdsaSigner(),
		},
		{
			name:   "Benchmark_Ed25519",
			signer: signers.NewEd25519Signer(),
		},
	}
	var err error
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			c := NewCertificate(tt.signer)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err = c.InitSelfSigned()
				if err != nil {
					b.Fatal("initSelfSigned")
				}
			}
			b.StopTimer()
		})
	}
}
