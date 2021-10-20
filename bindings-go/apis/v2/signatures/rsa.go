package signatures

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"

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

func (s RsaSigner) Sign(data []byte) (*v2.SignatureSpec, error) {
	signature, err := rsa.SignPKCS1v15(nil, &s.privateKey, crypto.SHA256, data)
	if err != nil {
		return nil, fmt.Errorf("failed signing hash, %w", err)
	}
	return &v2.SignatureSpec{
		Algorithm: "RSASSA-PKCS1-V1_5-SIGN", //TODO: check
		Data:      hex.EncodeToString(signature),
	}, nil
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
	key := untypedKey.(*rsa.PublicKey)
	return &RsaVerifier{
		publicKey: *key,
	}, nil
}

func (v RsaVerifier) Verify(signature v2.Signature) error {
	decodedHash, err := hex.DecodeString(signature.Digest.Value)
	if err != nil {
		return fmt.Errorf("failed decoding hash %s: %w", signature.Digest.Value, err)
	}
	decodedSignature, err := hex.DecodeString(signature.Signature.Data)
	if err != nil {
		return fmt.Errorf("failed decoding hash %s: %w", signature.Digest.Value, err)
	}
	err = rsa.VerifyPKCS1v15(&v.publicKey, crypto.SHA256, decodedHash, decodedSignature)
	if err != nil {
		return fmt.Errorf("signature verification failed, %w", err)
	}
	return nil
}
