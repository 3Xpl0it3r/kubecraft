package certs

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
)

func TryCreateCACertAndKeyFiles() {
	if certOrKeyExist(DefaultPkiPath) {
		return
	}
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	config := certutil.Config{
		CommonName:   DefaultCommonName,
		Organization: []string{},
		AltNames:     certutil.AltNames{},
		Usages:       []x509.ExtKeyUsage{},
		NotBefore:    expirationTime,
	}

	key, _ := NewPrivateKey(EncryptionAlgorithmRSA2048)
	ca, _ := NewSelfSignedCACert(config.CommonName, nil, key)

	if err := writeKey(pathForKey(DefaultPkiPath, "ca"), key); err != nil {
		panic("write key failed")
	}
	if err := writeCert(pathForCert(DefaultPkiPath, "ca"), ca); err != nil {
		panic("write ca certfaield")
	}

}

func certOrKeyExist(pkiPath string) bool {

	_, certErr := os.Stat(pathForCert(DefaultPkiPath, "ca"))
	_, keyErr := os.Stat(pathForKey(DefaultPkiPath, "ca"))

	return !(os.IsNotExist(certErr) && os.IsNotExist(keyErr))
}

func writeKey(privatePath string, key crypto.Signer) error {
	encoded, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		return errors.Wrapf(err, "unable to unmarshal private key to PEM")
	}
	if err := keyutil.WriteKey(privatePath, encoded); err != nil {
		return errors.Wrapf(err, "unable to unmarshal private key to PEM")
	}
	return nil
}

func writeCert(certPath string, cert *x509.Certificate) error {
	if err := certutil.WriteCert(certPath, pem.EncodeToMemory(&pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	})); err != nil {
		return errors.Wrapf(err, "unable to write cacert to pem")
	}
	return nil
}
