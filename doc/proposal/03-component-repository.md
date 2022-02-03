# Component Repository Specification

*Component Descriptors* are stored in *Component Repositories*. This chapter describes in a programing language/protocol
agnostic fashion the functions a *Component Repository* must implement to be OCM conform.

Usually sources and resources are stored in some external storage, e.g. a docker image is stored in some OCI registry
and the *Component Descriptor* contains only the information how to access it. As an alternative *Component Repositories* 
also provide the possibility to store technical artefacts together with the *Component Descriptors* in the 
*Component Repository* itself as so-called *local blobs*. This allows to pack all component versions with their 
technical artefacts in a *Component Repository* as a completely self-contained package. This is a typical requirement 
if you need to deliver your product into a fenced landscape. 

Local blobs could be referenced in the *resources* section of a *Component Descriptor, e.g. if the blob contains a 
binary, or in the *sources* section, e.g. if the blob contains a tar archive of a git repo. 

As we assume that OCI repositories will become the leading storage technology for technical artefacts and will
often be the backend of a *Component Repository* implementations, *Component Repositories* provide methods
to store OCI artefacts as so-called local OCI artefacts. It is not required that these OCI artefacts are accessible by 
an external OCI HTTP endpoint as specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md). 
If the repository provides such an HTTP endpoint, it could be requested by the method GetOciEndpointForLocalOciArtefact. 

The idea behind requiring methods for local OCI artefacts is the possibility is to define a file based transport format 
and implement a *Component Repository* interface on top of it. This allows implementing e.g. transport algorithms just 
working with the *Component Repository* interface on the source *Component Repository*, the transport format and the 
target *Component Repository* with no need to convert OCI artefacts into local blobs and vice versa. 

Local OCI blobs could only be referenced in the *resources* section of a *Component Descriptor.

Particular *Component Repository* implementations might extend the interface by special methods for other
types of binaries, e.g. helm charts.

**Todo: Image mit CDs, localblobs, localociimages, ...**

## Functions for Component Descriptors

### UploadComponentDescriptor

**Description**: Uploads a *Component Descriptor* to the *Component Repository*. If successful the *Component Descriptor*
is accessible by its name and version (see GetComponentDescriptor). The name and version of a *Component Descriptor*
is the identifier of a *Component Descriptor*, therefore if there already exists a *Component Descriptor*
with the same name and version, the upload fails. 

A *Component Repository* must check if all referenced *Component Descriptors*, local blobs and local OCI artefacts 
are already stored in the *Component Repository*.

If the identifier of entries in *resources*, *sources* or *componentReferences* are not unique, as described before,
an *invalidArgument* error must be returned.

If the identifier of *resources* and *sources* (name plus extraIdentity) are not unique, a *Component Repository* might 
automatically add the version field to the extraIdentity to resolve this problem. Of course, it must still fail, if 
uniqueness of the resource identifier is not achieved this way.

If the last entry in the "repositoryContext" field, is not an entry for the current *Component Repository*, such an 
entry is automatically added.

**Inputs**:

- String componentDescriptor: Json or Yaml representation of the *Component Descriptor*

**Outputs**:

**Errors**:

- alreadyExists: If there already exists a *Component Descriptor* with the same name and version
- missingReference: If a referenced *Component Descriptor*, local blob or local OCI artefact does not exist in 
  the *Component Repository*
- invalidArgument: If the parameter *componentDescriptor* is missing or is not conform to the
    [specified json schema](component-descriptor-v2-schema.yaml)
- repositoryError: If some error occurred in the *Component Repository*

### GetComponentDescriptor

**Description**: Returns the *Component Descriptor* as a Yaml string according the [JSON schema](component-descriptor-v2-schema.yaml).

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*

**Outputs**:

- String: Yaml string of the *Component Descriptor*

**Errors**:

- invalidArgument: If one of the input parameters is empty
- doesNotExist: If the *Component Descriptor* does not exist
- repositoryError: If some error occurred in the *Component Repository*

### DeleteComponentDescriptor

**Description**: Deletes the *Component Descriptor*. The deletion of a *Component Descriptor fails if it is referenced 
by another *Component Descriptor*.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*

**Outputs**:

**Errors**:

- doesNotExist: If the *Component Descriptor* does not exist
- existingReference: If the *Component Descriptor* is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*

## Functions for Local Blobs

### UploadLocalBlob

**Description**: Allows uploading binary data. The binary data belongs to a particular *Component Descriptor* 
and can be referenced by the component descriptor in its *resources* or *sources* section. 
*Component Descriptors* are not allowed to reference local blobs of other *Component Descriptors*. 

When uploading a local blob it is not required that the corresponding *Component Descriptor* already exists. 
Usually local blobs are uploaded first because it is not allowed to upload a *Component Descriptor* if its local 
blobs not already exist.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- BinaryStream data: Binary stream containing the local blob data.

**Outputs**:

- String blobIdentifier: The identifier, which must be used in the resource or source reference in the *Component Descriptor* 

**Errors**:

- invalidArgument: If one of the input parameters is empty or not valid
- repositoryError: If some error occurred in the *Component Repository*

**Example**: 
Assume you want to upload a helm chart, which is stored in a folder of your local file system. First you tar and compress
the folder. Then you upload it to your *Component Repository* with the *UploadLocalBlob* function and get the 
*blobIdentifier*: 

```
sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623
```

The entry in the *Component Descriptor* looks as follows:

```
...
  resources:
  - name: example-chart-blob
    type: helm.io/chart
    access:
      digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623
      type: localOciBlob
  ... 
```

### GetLocalBlob

**Description**: Fetches the binary data of a local blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String blobIdentifier: Identifier of the local blob

**Outputs**:

- BinaryStream data: Binary stream containing the local blob data.

**Errors**:

- doesNotExist: If the local blob does not exist
- invalidArgument: If one of the input parameters is empty or invalid
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalBlob

**Description**: Deletes a local blob. An error occurs if there is still an existing reference to the blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String blobIdentifier: Identifier of the local blob

**Outputs**:

**Errors**:

- doesNotExist: If the local blob does not exist
- existingReference: If the local blob is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*

## Functions for Local OCI Artefacts

OCI artefacts are usually referenced by more than one *Component Descriptor*. Therefore, local OCI artefacts are not 
stored in the context of a *Component Descriptor* as local blobs. Remember, when you upload a local blob, you need 
to specify the *Component Descriptor* it belongs to and only this *Component Descriptor* could reference it.

The semantics of the methods for handling OCI artefacts are quite similar to the 
[OCI specification](https://github.com/opencontainers/distribution-spec/blob/main/spec.md) but abstracts from the
protocol details.

### UploadLocalOciBlob

**Description**: Uploads an OCI blob to the *Component Repository*. The *name* and *tag* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the OCI artefact
- BinaryStream content: stream of the binary data of the blob

**Outputs**:

- String digest: digest to access the blob according to the OCI specification

**Errors**:

- invalidArgument: If one of the input parameters is empty or has the wrong format.
- repositoryError: If some error occurred in the *Component Repository*

### UploadLocalOCIManifest

**Description**: Uploads an OCI manifest. The manifest could be either an
[image manifest](https://github.com/opencontainers/image-spec/blob/main/manifest.md) or an
[image index manifest](https://github.com/opencontainers/image-spec/blob/main/image-index.md).

The *name* and *tag* of the OCI artefact must follow the rules specified 
[here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

The rules for rejecting a manifest are the same as specified in the 
[OCI distribution spec](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#push).

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the OCI artefact
- BinaryStream manifest: stream of the manifest

**Outputs**:

- String digest: digest to access the manifest according to the OCI specification

**Errors**:

- alreadyExists: If there exists a local OCI artefact with the same name and tag
- invalidArgument: If one of the input parameters is empty or has the wrong format.
- invalidArtefact: If the *Component Repository* found some inconsistencies with respect to the OCI manifest
- missingReference: If referenced blobs or manifests are missing
- repositoryError: If some error occurred in the *Component Repository*

### GetLocalOciBlob

**Description**: Returns an OCI blob of the *Component Repository*. The *name* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String digest: digest of the blob

**Outputs**:

- BinaryStream content: binary data of the blob

**Errors**:

- invalidArgument: If one of the input parameters is empty or has the wrong format.
- doesNotExist: If the blob does not exist
- repositoryError: If some error occurred in the *Component Repository*

### GetLocalOciManifest

**Description**: Returns manifest of the *Component Repository*. The *name* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

Either tag or manifest must be set, but not both.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the manifest
- String digest: digest of the manifest

**Outputs**:

- BinaryStream content: binary data of the manifest

**Errors**:

- invalidArgument: If one of the input parameters is empty or has the wrong format.
- doesNotExist: If the manifest does not exist
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalOciBlob

**Description**: Deletes an OCI blob of the *Component Repository*. The *name* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

It is up to the *Component Repository* if an error is thrown if the blob is still referenced by a manifest.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String digest: digest of the blob

**Outputs**:

- invalidArgument: If one of the input parameters is empty or has the wrong format.
- existingReference: vis still referenced by a manifest
- doesNotExist: if the blob does not exist
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalOciManifest

**Description**: Deletes an OCI manifest of the *Component Repository*. The *name* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

It is up to the *Component Repository* if an error is thrown if the manifest is still referenced by a manifest.

Either tag or manifest must be set, but not both.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the manifest
- String digest: digest of the manifest

**Outputs**:

- invalidArgument: If one of the input parameters is empty or has the wrong format.
- existingReference: if the blob is still referenced by a manifest
- doesNotExist: if the blob does not exist
- repositoryError: If some error occurred in the *Component Repository*

**Example**:
Assume you have uploaded a manifest for a local OCI artefact with name *example.com/test* and tag *1.1.1*, you could reference it 
in a *Component Descriptor* as follows:

```
...
  resources:
  - name: example-chart-blob
    type: ociImage
    access:
      ociArtefactName: example.com/test
      tag: 1.1.1
      type: localOciArtefact
  ... 
```

### GetOciEndpointForLocalOciArtefact

**Description**: If the *Component Repository* supports the [Open Container Initiative Distribution Specification]
(https://github.com/opencontainers/distribution-spec/blob/main/spec.md) and provides a spec conforming endpoint
for local OCI artefacts, it is possible that these are accessible by another external name, e.g. if you implement a solution
where several *Component Repositories* on top of only one OCI registry.

Providing the ociArtefactName, i.e. the name for the local OCI artefact, this method returns the external name. 

The external name needs not to be the same as *ociArtefactName*, but it is recommended 
that *ociArtefactName* is a suffix of the external name.

**Inputs**:

- String ociArtefactName: name of the OCI artefact

**Outputs**:

- String registryEndPoint: endpoint of the OCI registry
- String externalNam: external name of the OCI artefact

**Errors**:

- notSupported: If this kind of OCI access to the local OCI artefact is not supported by the *Component Repository*
- doesNotExist: If the local OCI artefact does not exist
- invalidArgument: If one of the input parameters is empty or has a wrong format
- repositoryError: If some error occurred in the *Component Repository*

**Example:**

Assume the method is called as follows

```
    GetOciEndpointForLocalOciArtefact("example.com/test")
```

and returns the following data

```
    registryEndPoint="test.example-registry.com
    name="examplefolder/example.com/test"
```

then the manifest could be fetched with an HTTP GET request using one of the following URLs

```
    test.example-registry.com/v2/examplefolder/example.com/test/manifests/sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623
```

or 

```
    test.example-registry.com/v2/examplefolder/example.com/test/manifests/0.1.1
```