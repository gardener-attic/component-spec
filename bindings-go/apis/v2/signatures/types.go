package signatures

import v2 "github.com/gardener/component-spec/bindings-go/apis/v2"

type Signer interface {
	Sign(data []byte) (*v2.SignatureSpec, error)
}

type Verifier interface {
	Verify(signature v2.Signature) error
}
