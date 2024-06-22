package certs

type EncryptionAlgorithmType string

const (
	DefaultCommonName   = "kubecraft"
	DefaultUserName     = "admin"
	DefaultOrganization = "kubecraft"

	OutDir                    = "pki"
	DefaultKubeConfigFileName = "kubeconfig"
	DefaultPkiPath            = OutDir

	DefaultApiServerCertName = "apiserver"
)

var (
	ApiServerCertFile = pathForCert(DefaultPkiPath, DefaultApiServerCertName)
	ApiServerKeyFile  = pathForKey(DefaultPkiPath, DefaultApiServerCertName)
)

const (
	// EncryptionAlgorithmECDSAP256 defines the ECDSA encryption algorithm type with curve P256.
	EncryptionAlgorithmECDSAP256 EncryptionAlgorithmType = "ECDSA-P256"
	// EncryptionAlgorithmRSA2048 defines the RSA encryption algorithm type with key size 2048 bits.
	EncryptionAlgorithmRSA2048 EncryptionAlgorithmType = "RSA-2048"
	// EncryptionAlgorithmRSA3072 defines the RSA encryption algorithm type with key size 3072 bits.
	EncryptionAlgorithmRSA3072 EncryptionAlgorithmType = "RSA-3072"
	// EncryptionAlgorithmRSA4096 defines the RSA encryption algorithm type with key size 4096 bits.
	EncryptionAlgorithmRSA4096 EncryptionAlgorithmType = "RSA-4096"
)
const (
	// PrivateKeyBlockType is a possible value for pem.Block.Type.
	PrivateKeyBlockType = "PRIVATE KEY"
	// PublicKeyBlockType is a possible value for pem.Block.Type.
	PublicKeyBlockType = "PUBLIC KEY"
	// CertificateBlockType is a possible value for pem.Block.Type.
	CertificateBlockType = "CERTIFICATE"
	// RSAPrivateKeyBlockType is a possible value for pem.Block.Type.
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
	// ECPrivateKeyBlockType is a possible value for pem.Block.Type.
	ECPrivateKeyBlockType = "EC PRIVATE KEY"
)
