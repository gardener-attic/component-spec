package signatures

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

type RsaSigner struct {
	privateKey rsa.PrivateKey
}

func CreateRsaSignerFromKeyFile(pathToPrivateKey string) (*RsaSigner, error) {
	privKeyFile, err := ioutil.ReadFile(pathToPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed opening private key file %w", err)
	}

	block, _ := pem.Decode([]byte(privKeyFile))
	if block == nil {
		return nil, fmt.Errorf("failed decoding PEM formatted block in key %w", err)
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing key %w", err)
	}
	return &RsaSigner{
		privateKey: *key,
	}, nil
}

func (s RsaSigner) Sign(componentDescriptor v2.ComponentDescriptor, digest v2.DigestSpec) (*v2.SignatureSpec, error) {
	decodedHash, err := hex.DecodeString(digest.Value)
	if err != nil {
		return nil, fmt.Errorf("failed decoding hash to bytes")
	}
	hashType, err := hashAlgorithmLookup(digest.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed looking up hash algorithm")
	}
	signature, err := rsa.SignPKCS1v15(nil, &s.privateKey, hashType, decodedHash)
	if err != nil {
		return nil, fmt.Errorf("failed signing hash, %w", err)
	}
	return &v2.SignatureSpec{
		Algorithm: "RSASSA-PKCS1-V1_5-SIGN",
		Data:      hex.EncodeToString(signature),
	}, nil
}

// maps a hashing algorithm string to crypto.Hash
func hashAlgorithmLookup(algorithm string) (crypto.Hash, error) {
	switch strings.ToUpper(algorithm) {
	case "SHA256":
		return crypto.SHA256, nil
	}
	return 0, fmt.Errorf("hash Algorithm %s not found", algorithm)
}

type RsaVerifier struct {
	publicKey rsa.PublicKey
}

func CreateRsaVerifierFromKeyFile(pathToPublicKey string) (*RsaVerifier, error) {
	publicKey, err := ioutil.ReadFile(pathToPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed opening public key file %w", err)
	}
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, fmt.Errorf("failed decoding PEM formatted block in key %w", err)
	}
	untypedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing key %w", err)
	}
	switch key := untypedKey.(type) {
	case *rsa.PublicKey:
		return &RsaVerifier{
			publicKey: *key,
		}, nil
	default:
		return nil, fmt.Errorf("public key format is not supported. Only rsa.PublicKey is supported")
	}
}

func (v RsaVerifier) Verify(componentDescriptor v2.ComponentDescriptor, signature v2.Signature) error {
	decodedHash, err := hex.DecodeString(signature.Digest.Value)
	if err != nil {
		return fmt.Errorf("failed decoding hash %s: %w", signature.Digest.Value, err)
	}
	decodedSignature, err := hex.DecodeString(signature.Signature.Data)
	if err != nil {
		return fmt.Errorf("failed decoding hash %s: %w", signature.Digest.Value, err)
	}
	algorithm, err := hashAlgorithmLookup(signature.Digest.Algorithm)
	if err != nil {
		return fmt.Errorf("failed looking up hash algorithm for %s: %w", signature.Digest.Algorithm, err)
	}
	err = rsa.VerifyPKCS1v15(&v.publicKey, algorithm, decodedHash, decodedSignature)
	if err != nil {
		return fmt.Errorf("signature verification failed, %w", err)
	}
	return nil
}
