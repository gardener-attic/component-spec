# Golang Binding for the ComponentDescriptor

:warning: currently only the ComponentDescriptor v2 is implemented

### Install

```
go get github.com/gardener/component-spec/bindings-go
```

### Usage

**Decode and Encode a Component Descriptor**

```go
package main

import (
    "fmt"

    "github.com/gardener/component-spec/bindings-go/codec"
    "github.com/gardener/component-spec/bindings-go/apis/v2"
    "github.com/gardener/component-spec/bindings-go/utils/selector"
)

func main() {
	data := []byte(`
meta:
  schemaVersion: 'v2'

component:
  name: 'github.com/gardener/gardener'
  version: 'v1.7.2'

  provider: internal

  repositoryContexts: []
  sources: []
  componentReferences: []

  resources:
  - name: 'apiserver'
    version: 'v1.7.2'
    type: 'ociImage'
    relation: local
    access:
      type: 'ociRegistry'
      imageReference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'

  - name: 'hyperkube'
    version: 'v1.16.4'
    type: 'ociImage'
    extraIdentity:
      myid: '1'
    relation: external
    access:
      type: 'ociRegistry'
      imageReference: 'k8s.gcr.io/hyperkube:v1.16.4'
  - name: 'hyperkube'
    version: 'v1.17.4'
    type: 'ociImage'
    extraIdentity:
      myid: '2'
    relation: external
    access:
      type: 'ociRegistry'
      imageReference: 'k8s.gcr.io/hyperkube:v1.16.4'
`)

	component := &v2.ComponentDescriptor{}
	err := codec.Decode(data, component)
	check(err)

	encData, err := codec.Encode(component)
	check(err)
    fmt.Println(string(encData)) // prints the components descriptor as json
}
```

##### Select Resources

:warning: Note that the following examples use the above described component descriptor.

**Decode access types**
```go
component := &v2.ComponentDescriptor{}
err := codec.Decode(data, component)
check(err)

// get a specific local resource
res, err := component.GetLocalResource(v2.OCIImageType, "apiserver", "v1.7.2")
check(err)
fmt.Printf("%#v\n", res)

// get a specific external resource
res, err = component.GetExternalResource(v2.OCIImageType, "hyperkube", "v1.16.4")
check(err)
fmt.Printf("%#v\n", res)

// get the access for a resource
// specific access type can be decoded using the access type codec.
accessTypeCodec := v2.NewCodec(nil, nil, nil)
ociAccess := &v2.OCIRegistryAccess{}
check(accessTypeCodec.Decode(res.Access.Raw, ociAccess))
fmt.Println(ociAccess.ImageReference) // prints: k8s.gcr.io/hyperkube:v1.16.4
```

**Select resources by identity**
```go
component := &v2.ComponentDescriptor{}
err := codec.Decode(data, component)
check(err)

// get a component by its identity via selectors
idSelector := selector.DefaultSelector{
    "name": "hyperkube",
}
resources, err := component.GetResourceBySelector(idSelector)
check(err)
fmt.Printf("%d\n", len(resources)) // prints "2" as both hyperkube images match the identity

// get a component by additional identity information
idSelector = selector.DefaultSelector{
    "name": "hyperkube",
    "myid": "1",
}
resources, err = component.GetResourceBySelector(idSelector)
check(err)
fmt.Printf("%d\n", len(resources)) // prints "1" as only one hyperkube image matches the name and myid attribute.
fmt.Printf("%s - %s\n", resources[0].Name, resources[0].Version) // prints "hyperkube - v1.16.4"
```

**Select resources by their identity using jsonschema**
```go
component := &v2.ComponentDescriptor{}
err := codec.Decode(data, component)
check(err)

schemaSelector, err := selector.NewJSONSchemaSelectorFromString(`
type: object
properties:
  name:
    type: string
    enum: ["hyperkube"]
  myid:
    type: string
    enum: ["1"]
`)
check(err)

resources, err = component.GetResourceBySelector(schemaSelector)
check(err)
fmt.Printf("%d\n", len(resources)) // prints "1" as only one hyperkube image matches the name and myid attribute.
fmt.Printf("%s - %s\n", resources[0].Name, resources[0].Version) // prints "hyperkube - v1.16.4"
```


For more examples see the [examples folder](./examples)
