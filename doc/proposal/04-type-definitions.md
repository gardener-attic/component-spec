# Type Definitions

This chapter defines the format of some types and their access format for resource and source references. It is 
recommended use these formats in all implementations based on OCM.

## Source and Resource References

### OCI Artefact

This is the recommended format to reference an OCI artefact in a 

```yaml
    name: example-name
    type: ociImage
    access:
      imageReference: name[:tag|@digest]
      type: ociRegistry
    
```

*imageReference* is the reference to the OCI image reference, whereby name, tag and digest are specified 
[here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull)

### Local Blob

```yaml
resources:
  - name: example-resource-local
    relation: local
    type: helm.io/chart
    version: v0.1.0
    access:
      digest: sha256:123...
      type: localOciBlob
```

### S3

```yaml
resources:
  - name: example-resource-s3
    relation: local
    type: example-type
    version: "v0.1.0"
    access:
      bucketName: examplebucket
      objectKey: objects/123...
      type: s3
```

