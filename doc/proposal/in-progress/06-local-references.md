# Local Blob Reference Definitions

## Local Blobs

**this is work in progress**

In a *Component Repository*, local blobs can be stored together with a *Component Descriptor*. Local blobs could store
different data like OCI images, helm charts, configuration data etc. 

```
...
  resources:
  - name: example-image
    type: oci-image
    access:
      type: localOciBlob
      mediaType: application/vnd.oci.image.manifest.v1+json
      annotations:
        name: test/monitoring
      localAccess: "digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
      globalAccess: 
        imageReference: somePrefix/test/monitoring@sha:...
        type: ociRegistry
... 
```

```yaml
resources:
  - name: example-name
    relation: local
    type: helm.io/chart
    version: v0.1.0
    access:
      digest: <identifier/digest of the local blob>
      type: localOciBlob
```


