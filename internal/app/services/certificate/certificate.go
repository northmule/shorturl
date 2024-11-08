package certificate

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

// Certificate сервис сертификата
type Certificate struct {
	keyPath    string
	certPath   string
	privateKey crypto.Signer
}

// NewCertificate конструктор
func NewCertificate() *Certificate {
	tmpPath := os.TempDir()
	return &Certificate{
		keyPath:  path.Join(tmpPath, "key.pem"),
		certPath: path.Join(tmpPath, "cert.pem"),
	}
}

// SetPrivateKey установка алгоритма шифрования
func (c *Certificate) SetPrivateKey(algo string) error {
	var err error
	var key crypto.Signer
	switch algo {
	case "rsa":
		key, err = rsa.GenerateKey(rand.Reader, 4096)

	case "ecdsa":
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	case "ed25519":
		_, key, err = ed25519.GenerateKey(rand.Reader)

	default:
		return errors.New("ожидаются: rsa | ecdsa | ed25519")

	}

	if err != nil {
		return err
	}
	c.privateKey = key

	return nil
}

// InitSelfSigned создаёт ключ и сертификат для TLS сервера
func (c *Certificate) InitSelfSigned() error {
	var err error

	if c.privateKey == nil {
		return errors.New("first you need to call SetPrivateKey(algo)")
	}

	privateKey := c.privateKey
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	// серийный номер сертификата
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	// шаблон сертификата
	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Horns and Hooves"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost"},
		NotBefore:   time.Now(),
		// время жизни сертификата
		NotAfter: time.Now().Add(24 * time.Hour),

		KeyUsage: x509.KeyUsageDigitalSignature,
		// Набобр для вариантов использования, example: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// сертификат
	certBytes, err := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, privateKey.Public(), privateKey)
	if err != nil {
		return err
	}
	// кодирование сертификата в pem для хранения
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if pemCert == nil {
		return err
	}

	// Запись сертификата в файл
	if err = os.WriteFile(c.certPath, pemCert, 0644); err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	// кодирование ключа в pem для хранения
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})
	if pemKey == nil {
		return errors.New("failed to encode key to PEM")
	}
	// Запись ключа в файл
	if err = os.WriteFile(c.keyPath, pemKey, 0600); err != nil {
		return err
	}

	return nil
}

// CertPath путь к сертфификату
func (c *Certificate) CertPath() string {
	return c.certPath
}

// KeyPath путь к ключу
func (c *Certificate) KeyPath() string {
	return c.keyPath
}
