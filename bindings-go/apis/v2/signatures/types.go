package signatures

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"strings"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

type Signer interface {
	// Sign returns the signature for the data for the component-descriptor
	Sign(componentDescriptor v2.ComponentDescriptor, digest v2.DigestSpec) (*v2.SignatureSpec, error)
}

type Verifier interface {
	// Verify checks the signature, returns an error on verification failure
	Verify(componentDescriptor v2.ComponentDescriptor, signature v2.Signature) error
}

type Hasher struct {
	HashFunction  hash.Hash
	AlgorithmName string
}

func HasherForName(algorithmName string) (*Hasher, error) {
	switch strings.ToUpper(algorithmName) {
	case "SHA256":
		return &Hasher{
			HashFunction:  sha256.New(),
			AlgorithmName: "SHA256",
		}, nil
	}
	return nil, fmt.Errorf("hash algorithm %s not found/implemented", algorithmName)
}
