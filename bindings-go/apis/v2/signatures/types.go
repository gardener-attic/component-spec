package signatures

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"strings"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// Signer interface is used to implement different signing algorithms.
// Each Signer should have a matching Verifier.
type Signer interface {
	// Sign returns the signature for the data for the component-descriptor
	Sign(componentDescriptor cdv2.ComponentDescriptor, digest cdv2.DigestSpec) (*cdv2.SignatureSpec, error)
}

// Verifier interface is used to implement different verification algorithms.
// Each Verifier should have a matching Signer.
type Verifier interface {
	// Verify checks the signature, returns an error on verification failure
	Verify(componentDescriptor cdv2.ComponentDescriptor, signature cdv2.Signature) error
}

// Hasher encapsulates a hash.Hash interface with an algorithm name.
type Hasher struct {
	HashFunction  hash.Hash
	AlgorithmName string
}

const SHA256 = "sha256"

// HasherForName creates a Hasher instance for the algorithmName.
func HasherForName(algorithmName string) (*Hasher, error) {
	switch strings.ToLower(algorithmName) {
	case SHA256:
		return &Hasher{
			HashFunction:  sha256.New(),
			AlgorithmName: SHA256,
		}, nil
	case strings.ToLower(cdv2.NoDigest):
		return &Hasher{
			HashFunction:  nil,
			AlgorithmName: cdv2.NoDigest,
		}, nil
	}
	return nil, fmt.Errorf("hash algorithm %s not found/implemented", algorithmName)
}

type ResourceDigester interface {
	DigestForResource(ctx context.Context, componentDescriptor cdv2.ComponentDescriptor, resource cdv2.Resource, hasher Hasher) (*cdv2.DigestSpec, error)
}
