# Type Definitions

A component descriptor may describe references to sources and resources. For every source and resource it must be 
defined how to access it. The format of the access data can depend on the type of a resource and the
repository where it is stored. This chapter defines some formats of access data. It is 
recommended use these formats in all implementations based on OCM.

## Source and Resource References

### OCI Artefact

This is the recommended format to reference an OCI artefact in an OCI image registry:

```yaml
resources:
  - name: example-name
    type: ociImage
    access:
      imageReference: name[:tag|@digest]
      type: ociRegistry
```

*imageReference* is the reference to the OCI image reference, whereby name, tag and digest are specified 
[here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull)

### Local Blob

In a component repository, local blobs can be stored together with a component descriptor. Every local blob must have
a blob identifier. In the case of an [OCI component repository](05-component-repository-oci.md), this is the digest 
of the blob.

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

### S3

The example below shows the recommended format to define the access of an object in an s3 bucket.

```yaml
resources:
  - name: example-name
    relation: local
    type: example-type
    version: "v0.1.0"
    access:
      bucketName: <s3 bucket name>
      objectKey: <s3 object key>
      type: s3
```
