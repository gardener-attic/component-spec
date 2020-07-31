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

    "github.com/gardener/component-spec/client-go" 
    "github.com/gardener/component-spec/client-go/apis/v2"
)

func main() {
    data := []byte(`
meta:
  schemaVersion: v2

components:
- name: my-comp
  version: 1.2.3
  type: gardenerComponent

  dependencies:
  - name: apiserver
    version: 1.16.4
    type: ociImage
    imageReference: eu.gcr.io/gardener-project/hyperkube:1.16.4
`)
    cd := &v2.ComponentDescriptor{}
    err := client_go.Decode(data, cd)
    check(err)
    
    // apply overrides 
    // and create a new component descriptor with overwritten dependencies
    effectiveCD, err := cd.ApplyOverwrites()
    check(err)

    // get a specific component
    myComp, err := effectiveCD.Components.GetComponentByMetadata(v2.ComponentMetadata{Name: "my-comp", Version: "1.2.3", Type: v2.GardenerComponentType})
    check(err)

    // get a dependency of a component
    apiserverDep, err := myComp.GetDependencies().GetComponentByMetadata(v2.ComponentMetadata{Name: "apiserver", Version: "1.16.4", Type: v2.OCIComponentType})
    check(err)
    // known types implement the ComponentAccessor interface and can be cast to the specific type.
    apiserverImg := apiserverDep.(*v2.OCIImage)
    
    fmt.Println(apiserverImg.ImageReference) // stdout: eu.gcr.io/gardener-project/hyperkube:1.16.4
}
```