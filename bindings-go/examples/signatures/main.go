package main

import (
	"context"
	"flag"
	"fmt"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

var privateKeyPath *string
var publicKeyPath *string

func init() {
	privateKeyPath = flag.String("private-key", "private", "private key for signing")
	publicKeyPath = flag.String("public-key", "public", "public key for verification")
}

func main() {
	flag.Parse()

	resAccess, err := cdv2.NewUnstructured(cdv2.NewGitHubAccess("url2", "ref", "commit"))
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	cd := cdv2.ComponentDescriptor{
		Metadata: cdv2.Metadata{
			Version: "v2",
		},
		ComponentSpec: cdv2.ComponentSpec{
			ObjectMeta: cdv2.ObjectMeta{
				Name:    "CD-Name<html>cool</html> Unicode ♥ unprintable characters \u0007 \u0031",
				Version: "v0.0.1",
			},
			ComponentReferences: []cdv2.ComponentReference{
				{
					Name:          "compRefName",
					ComponentName: "compRefNameComponentName",
					Version:       "v0.0.2compRef",
					ExtraIdentity: cdv2.Identity{
						"refKey": "refName",
					},
					Digest: &cdv2.DigestSpec{
						HashAlgorithm:          signatures.SHA256,
						NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
						Value:                  "value",
					},
				},
			},
			Resources: []cdv2.Resource{
				{
					IdentityObjectMeta: cdv2.IdentityObjectMeta{
						Name:    "Resource1",
						Version: "v0.0.3resource",
						ExtraIdentity: cdv2.Identity{
							"key": "value",
						},
					},
					Access: &resAccess,
					Digest: &cdv2.DigestSpec{
						HashAlgorithm:          signatures.SHA256,
						NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
						Value:                  "value",
					},
				},
			},
		},
	}
	ctx := context.TODO()
	err = signatures.AddDigestsToComponentDescriptor(ctx, &cd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
		return &cdv2.DigestSpec{
			HashAlgorithm:          signatures.SHA256,
			NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
			Value:                  "value",
		}, nil
	}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
		return &cdv2.DigestSpec{
			HashAlgorithm:          signatures.SHA256,
			NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
			Value:                  "value",
		}, nil
	})
	if err != nil {
		fmt.Printf("ERROR addingDigestsToComponentDescriptor %s", err)
		return
	}

	hasher, err := signatures.HasherForName(signatures.SHA256)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	norm, err := signatures.HashForComponentDescriptor(cd, *hasher)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	fmt.Println(norm.Value)

	signer, err := signatures.CreateRSASignerFromKeyFile(*privateKeyPath, cdv2.MediaTypePEM)
	if err != nil {
		fmt.Printf("ERROR create signer: %s", err)
		return
	}

	err = signatures.SignComponentDescriptor(&cd, signer, *hasher, "mySignatureName")
	if err != nil {
		fmt.Printf("ERROR sign: %s", err)
		return
	}

	verifier, err := signatures.CreateRSAVerifierFromKeyFile(*publicKeyPath)
	if err != nil {
		fmt.Printf("ERROR create verifier: %s", err)
		return
	}
	err = signatures.VerifySignedComponentDescriptor(&cd, verifier, "mySignatureName")
	if err != nil {
		fmt.Printf("ERROR verify signature: %s", err)
		return
	}
	fmt.Println("If not error is printed, successful")
}
