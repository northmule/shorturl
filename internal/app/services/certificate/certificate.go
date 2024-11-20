package certificate

import (
	"crypto"
	"crypto/rand"
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
	generator KeyGenerator
	keyPath   string
	certPath  string
}

// KeyGenerator интерфейс получения получение crypto.Signer
type KeyGenerator interface {
	GenerateKey() (crypto.Signer, error)
}

// NewCertificate конструктор
func NewCertificate(generator KeyGenerator) *Certificate {
	tmpPath := os.TempDir()
	return &Certificate{
		generator: generator,
		keyPath:   path.Join(tmpPath, "key.pem"),
		certPath:  path.Join(tmpPath, "cert.pem"),
	}
}

// InitSelfSigned создаёт ключ и сертификат для TLS сервера
func (c *Certificate) InitSelfSigned() error {
	var err error
	privateKey, err := c.generator.GenerateKey()
	if err != nil {
		return err
	}
	if privateKey == nil {
		return errors.New("GenerateKey empty")
	}
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
