// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
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
  - name: 'ftpRes'
    version: 'v1.7.2'
    type: 'custom1'
    access:
      type: 'x-ftp'
      url: ftp://example.com/my-resource

  externalResources:
  - name: 'nodeMod'
    version: '0.0.1'
    type: 'nodeModule'
    access:
      type: 'npm'
      nodeModule: my-module
      version: 0.0.1
`)
	// register additional types
	v2.KnownAccessTypes[NPMType] = npmCodec

	component := &v2.ComponentDescriptor{}
	err := codec.Decode(data, component)
	check(err)

	res, err := component.GetLocalResource("custom1", "ftpRes", "v1.7.2")
	check(err)
	// unknown types are serialized as custom type
	ftpAccess := res.Access.(*v2.CustomAccess)
	fmt.Println(ftpAccess.Data["url"]) // prints: ftp://example.com/my-resource

	res, err = component.GetExternalResource("nodeModule", "nodeMod", "0.0.1")
	check(err)
	npmAccess := res.Access.(*NPMAccess)
	fmt.Println(npmAccess.NodeModule) // prints: my-module
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

const NPMType = "npm"

type NPMAccess struct {
	v2.ObjectType
	NodeModule string `json:"nodeModule"`
	Version    string `json:"version"`
}

var _ v2.AccessAccessor = &NPMAccess{}

func (n NPMAccess) GetData() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NPMAccess) SetData(bytes []byte) error {
	var newNPM NPMAccess
	if err := json.Unmarshal(bytes, &newNPM); err != nil {
		return err
	}

	n.NodeModule = newNPM.NodeModule
	n.Version = newNPM.Version
	return nil
}

var npmCodec = &v2.AccessCodecWrapper{
	AccessDecoder: v2.AccessDecoderFunc(func(data []byte) (v2.AccessAccessor, error) {
		var npm NPMAccess
		if err := json.Unmarshal(data, &npm); err != nil {
			return nil, err
		}
		return &npm, nil
	}),
	AccessEncoder: v2.AccessEncoderFunc(func(accessor v2.AccessAccessor) ([]byte, error) {
		npm, ok := accessor.(*NPMAccess)
		if !ok {
			return nil, fmt.Errorf("accessor is not of type %s", NPMType)
		}
		return json.Marshal(npm)
	}),
}
