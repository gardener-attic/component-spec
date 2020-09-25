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

components:
- component:
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
- component:
    name: 'github.com/gardener/etcd'
    version: 'v1.3.0'

    provider: internal

    sources: []
    references: []
    localResources: []

    externalResources:
    - name: 'etcd'
      version: 'v3.5.4'
      type: 'ociImage'
      access:
        type: 'ociRegistry'
        imageReference: 'quay.io/coreos/etcd:v3.5.4'
`)

	list := &v2.ComponentDescriptorList{}
	err := codec.Decode(data, list)
	check(err)

	// get component by name and version
	comp, err := list.GetComponent("github.com/gardener/etcd", "v1.3.0")
	check(err)

	fmt.Println(comp.ExternalResources[0].Name) // prints: etcd

	// get a component by its name
	// The method returns a list as there could be multiple components with the same name but different version
	comps := list.GetComponentByName("github.com/gardener/gardener")
	fmt.Println(len(comps)) // prints: 1
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
