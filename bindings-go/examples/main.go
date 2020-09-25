// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/codec"
)

func main() {
	data := []byte(`
meta:
  schemaVersion: 'v2'

component:
  name: 'github.com/gardener/gardener'
  version: 'v1.7.2'

  provider: internal

  sources: []
  references: []

  localResources:
  - name: 'apiserver'
    version: 'v1.7.2'
    type: 'ociImage'
    access:
      type: 'ociRegistry'
      imageReference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'

  externalResources:
  - name: 'hyperkube'
    version: 'v1.16.4'
    type: 'ociImage'
    access:
      type: 'ociRegistry'
      imageReference: 'k8s.gcr.io/hyperkube:v1.16.4'
`)

	component := &v2.ComponentDescriptor{}
	err := codec.Decode(data, component)
	check(err)

	// get a specific local resource
	res, err := component.GetLocalResource(v2.OCIImageType, "apiserver", "v1.7.2")
	check(err)
	fmt.Printf("%v\n", res)

	// get a specific external resource
	res, err = component.GetExternalResource(v2.OCIImageType, "hyperkube", "v1.16.4")
	check(err)
	fmt.Printf("%v\n", res)

	// get the access for a resource
	// known types implement the AccessAccessor interface and can be cast to the specific type.
	ociAccess := res.Access.(*v2.OCIRegistryAccess)
	fmt.Println(ociAccess.ImageReference) // prints: eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
