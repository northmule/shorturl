package certificate

import (
	"os"
	"testing"
)

func TestInitSelfSigned_rsa(t *testing.T) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("rsa")
	if err != nil {
		t.Fatalf("SetPrivateKey(rsa) error = %v", err)
	}

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
}

func TestInitSelfSigned_ecdsa(t *testing.T) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("ecdsa")
	if err != nil {
		t.Fatalf("SetPrivateKey(ecdsa) error = %v", err)
	}

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
}

func TestInitSelfSigned_ed25519(t *testing.T) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("ed25519")
	if err != nil {
		t.Fatalf("SetPrivateKey(ecdsa) error = %v", err)
	}

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
}

func TestInitSelfSigned_InvalidAlgo(t *testing.T) {
	c := NewCertificate()
	err := c.SetPrivateKey("invalid")
	if err == nil {
		t.Fatalf("Expected error for invalid algo, but got nil")
	}

	err = c.InitSelfSigned()
	if err == nil {
		t.Fatalf("Expected error for invalid algo, but got nil")
	}
}

func BenchmarkCertificate_InitSelfSigned_rsa(b *testing.B) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("rsa")
	if err != nil {
		b.Fatal("setPrivateKey error")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = c.InitSelfSigned()
		if err != nil {
			b.Fatal("initSelfSigned")
		}
	}

}

func BenchmarkCertificate_InitSelfSigned_ecdsa(b *testing.B) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("ecdsa")
	if err != nil {
		b.Fatal("setPrivateKey error")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = c.InitSelfSigned()
		if err != nil {
			b.Fatal("initSelfSigned")
		}
	}

}

func BenchmarkCertificate_InitSelfSigned_ed25519(b *testing.B) {
	var err error
	c := NewCertificate()
	err = c.SetPrivateKey("ed25519")
	if err != nil {
		b.Fatal("setPrivateKey error")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = c.InitSelfSigned()
		if err != nil {
			b.Fatal("initSelfSigned")
		}
	}

}
