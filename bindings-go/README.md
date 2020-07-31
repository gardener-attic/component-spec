# Golang Binding for the ComponentDescriptor

:warning: currently only the ComponentDescriptor v2 is implemented

### Install

```
go get github.com/gardener/component-spec/client-go
```

### Usage

```go
package main

import (
    "fmt"

    "github.com/gardener/component-spec/bindings-go/codec" 
    "github.com/gardener/component-spec/bindings-go/apis/v2"
)

func main() {
    data := []byte(`
meta:
  schemaVersion: 'v2'

component:
  name: 'github.com/gardener/gardener'
  version: 'v1.7.2'

  externalResources:
  - name: 'hyperkube'
    version: 'v1.16.4'
    type: 'ociImage'
    access:
      type: 'ociRegistry'
      imageReference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'
`)
    component := &v2.ComponentDescriptor{}
    err := codec.Decode(data, component)
    check(err)

    // get a specific local resource
    res, err := component.GetLocalResource(v2.OCIImageType, "apiserver", "v1.7.2")
    check(err)
    fmt.Printf("%v", res)

    // get a specific external resource
    res, err = component.GetExternalResource(v2.OCIImageType, "hyperkube", "v1.16.4")
    check(err)
    fmt.Printf("%v", res)

    // get the access for a resource
    // known types implement the AccessAccessor interface and can be cast to the specific type.
    ociAccess := res.Access.(*v2.OCIRegistryAccess)
    fmt.Println(ociAccess.ImageReference) // prints: eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2
}
```

For more examples see the [examples folder](./examples)