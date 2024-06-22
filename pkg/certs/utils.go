package certs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"path/filepath"
	"time"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"

	"github.com/pkg/errors"
)

func NewPrivateKey(keyType EncryptionAlgorithmType) (crypto.Signer, error) {
	if keyType == EncryptionAlgorithmECDSAP256 {
		return ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	}

	rsaKeySize := rsaKeySizeFromAlgorithmType(keyType)
	if rsaKeySize == 0 {
		return nil, errors.Errorf("cannot obtain key size from unknown RSA algorithm: %q", keyType)
	}
	return rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
}

// copyd from kubeadm
func NewSignedCert(commandName string, organization []string, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer, isCA bool) (*x509.Certificate, error) {
	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64-1))

	if err != nil {
		return nil, err
	}

	serial = new(big.Int).Add(serial, big.NewInt(1))
	if len(commandName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}
	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	if isCA {
		keyUsage |= x509.KeyUsageCertSign
	}
	notBefore := caCert.NotBefore
	notAfter := caCert.NotAfter
	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   commandName,
			Organization: organization,
		},
		DNSNames:              []string{},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		SerialNumber:          serial,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &certTmpl, caCert, key.Public(), caKey)

	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func NewSelfSignedCACert(commonName string, organization []string, key crypto.Signer) (*x509.Certificate, error) {
	// returns a uniform random value in [0, max-1), then add 1 to serial to make it a uniform random value in [1, max).
	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64-1))
	if err != nil {
		return nil, err
	}
	serial = new(big.Int).Add(serial, big.NewInt(1))

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign

	notBefore := time.Now().UTC()

	notAfter := notBefore.Add(365 * 24 * time.Hour)

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: organization,
		},
		DNSNames:              []string{commonName},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func NewCertAndKey(caCert *x509.Certificate, caKey crypto.Signer, keyType EncryptionAlgorithmType) (*x509.Certificate, crypto.Signer, error) {
	key, err := NewPrivateKey(EncryptionAlgorithmRSA2048)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create private key")
	}

	cert, err := NewSignedCert(DefaultCommonName, nil, key, caCert, caKey, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to sign certificate")
	}
	return cert, key, nil
}

func TryLoadCertAndKeyFromDisk(pkiPath, name string) (*x509.Certificate, crypto.Signer, error) {
	certificatePath := pathForCert(pkiPath, name)
	certs, err := certutil.CertsFromFile(certificatePath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load the certificate file %s", certificatePath)
	}

	privateKeyPath := pathForKey(pkiPath, name)
	privKey, err := keyutil.PrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load the private key file %s", privateKeyPath)
	}

	var key crypto.Signer
	switch k := privKey.(type) {
	case *rsa.PrivateKey:
		key = k
	case *ecdsa.PrivateKey:
		key = k
	default:
		return nil, nil, errors.Errorf("the private key file %s is neither in RSA nor ECDSA format", privateKeyPath)
	}

	return certs[0], key, nil
}

func CreateCertAndKeyFilesWithCA(caCert *x509.Certificate, caKey crypto.Signer) {
	cert, key, err := NewCertAndKey(caCert, caKey, EncryptionAlgorithmRSA2048)
	if err != nil {
		panic(fmt.Sprintf("Generate ApiServer Cert and Key failed %w", err))
	}

	if err := writeKey(pathForKey(DefaultPkiPath, DefaultApiServerCertName), key); err != nil {
		panic(err)
	}
	if err := writeCert(pathForCert(DefaultPkiPath, DefaultApiServerCertName), cert); err != nil {
		panic(err)
	}

}

func TryLoadCertChainFromDisk(pkiPath, name string) (*x509.Certificate, []*x509.Certificate, error) {
	certificatePath := pathForCert(pkiPath, name)

	certs, err := certutil.CertsFromFile(certificatePath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load the certificate file %s", certificatePath)
	}

	cert := certs[0]
	intermediates := certs[1:]

	return cert, intermediates, nil
}

func rsaKeySizeFromAlgorithmType(keyType EncryptionAlgorithmType) int {
	switch keyType {
	case EncryptionAlgorithmRSA2048, "":
		return 2048
	case EncryptionAlgorithmRSA3072:
		return 3072
	case EncryptionAlgorithmRSA4096:
		return 4096
	default:
		return 0
	}
}

func pathForCert(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.crt", name))
}

func pathForKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.key", name))
}

func pathForPublicKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.pub", name))
}

func pathForCSR(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.csr", name))
}

func encodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

func encodPrivateKeyToPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	switch t := privateKey.(type) {
	case *ecdsa.PrivateKey:
		derBytes, err := x509.MarshalECPrivateKey(t)
		if err != nil {
			return nil, err
		}
		block := &pem.Block{
			Type:  ECPrivateKeyBlockType,
			Bytes: derBytes,
		}
		return pem.EncodeToMemory(block), nil
	case *rsa.PrivateKey:
		block := &pem.Block{
			Type:  RSAPrivateKeyBlockType,
			Bytes: x509.MarshalPKCS1PrivateKey(t),
		}
		return pem.EncodeToMemory(block), nil
	default:
		return nil, fmt.Errorf("private key is not a recognized type: %T", privateKey)
	}
}
