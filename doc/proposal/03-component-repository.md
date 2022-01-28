# Component Repository Specification

*Component Descriptors* are stored in *Component Repositories*. This chapter describes in a programing language/protocol
agnostic fashion the functions a *Component Repository* must implement to be OCM conform.

To allow storing all technical artefacts together with the *Component Descriptors* of a software product version,
these could also be stored in a *Component Repository* as so-called local blobs.

As we assume that OCI repositories will become the leading storage technology for technical artefacts and will
often be the backend of *Component Repository* Implementations, *Component Repositories* provide methods
to store OCI artefacts as so-called local OCI artefacts, i.e. OCI artefacts which are accessible by a special function
*getLocalOCIArtefact*. It is not required that these OCI artefacts are accessible by an external OCI HTTP endpoint as 
specified [here](https://docs.docker.com/registry/spec/api/). If the repository provides such an HTTP endpoint, it could 
be requested by the method getOciEndpointForLocalOciArtefact. This allows to e.g. transport images from one *Component
Repository* to another without the need to transform them into local blobs first. 

Particular *Component Repository* implementations might extend the interface by special methods handling other
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

If the last entry in the 


**Inputs**:

- String componentDescriptor: Json or Yaml representation of the *Component Descriptor*

**Outputs**:

**Errors**:

- alreadyExists: If there already exists a *Component Descriptor*
  with the same name and version
- missingReference: If a referenced *Component Descriptor*, local blob or local OCI artefact does not exist in 
  the *Component Repository*
- invalidArgument: If the parameter *componentDescriptor* is missing or is not conform to the
    [specified json schema](component-descriptor-v2-schema.yaml)
- repositoryError: If some error occurred in the *Component Repository*

### GetComponentDescriptor

**Description**: Returns the *Component Descriptor* as a Yaml string.

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

- existingReference: If the *Component Descriptor* is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*

## Functions for Local Blobs

### UploadLocalBlob

**Description**: Allows uploading binary data. The binary data belongs to a particular *Component Descriptor* 
and can be referenced by the component descriptor in its *resources* section. *Component Descriptors* are not allowed
to reference local blobs of other *Component Descriptors*. 

When uploading a local blob it is not required that the corresponding *Component Descriptor* already exists. 
Usually local blobs are uploaded first because it is not allowed to upload a *Component Descriptor* if its local 
blobs not already exists.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- BinaryStream data: Binary stream containing the local blob data.

**Outputs**:

- String blobIdentifier: The identifier, which must be used in the resource reference in the *Component Descriptor* 

**Errors**:

- invalidArgument: If one of the input parameters is empty
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
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalBlob

**Description**: Deletes a local blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String blobIdentifier: Identifier of the local blob

**Outputs**:

**Errors**:

- existingReference: If the local blob is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*

## Functions for Local OCI Artefacts

OCI artefacts and here mainly OCI images are usually referenced by more than one *Component Descriptor*. Therefore
we decided not to store local OCI artefacts in the context of a *Component Descriptor* as local blobs. Remember, when you 
upload a local blob you need to specify the *Component Descriptor* it belongs to. Only this *Component Descriptor* could 
reference it.

### UploadLocalOCIArtefact

**Description**: Uploads an OCI artefact to the *Component Repository*. The *name* and *tag* of the OCI artefact
must follow the rules specified [here](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull).

The details how the layer information is provided, e.g. as a tar with an entry for every layer, are implementation specific.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: reference of the OCI artefact
- String artefactMediaType: media type of the whole oci artefact
- String configMediaType: media type of the config of the OCI artefact
- Array[String] layerMediaTypes: the list of media types for each layer, whereby the ith entry belongs to the ith layer 
- BinaryStream config: stream of the data of the config blob
- BinaryStream layer: stream of the binary data of the different layer

**Outputs**:

**Errors**:

- alreadyExists: If there exists a local OCI artefact with the same name and tag
- invalidArgument: If one of the input parameters is empty or has the wrong format
- repositoryError: If some error occurred in the *Component Repository*

**Example**:
Assume you have uploaded a local OCI artefact with name *example.com/test* and tag *1.1.1* you could reference it 
in a *Component Descriptor* as follows:

```
...
  resources:
  - name: example-chart-blob
    type: ociImage
    access:
      imageReference: example.com/test
      type: localOciArtefact
  ... 
```

### GetLocalOCIArtefact

**Description**: Returns the local OCI artefact. The details how the layer information is provided, e.g. as a tar with 
an entry for every layer, are implementation specific.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: reference of the OCI artefact

**Outputs**:

- String artefactMediaType: media type of the whole oci artefact
- String configMediaType: media type of the config of the OCI artefact
- Array[String] layerMediaTypes: the list of media types for each layer, whereby the ith entry belongs to the ith layer
- BinaryStream config: stream of the data of the config blob
- BinaryStream layer: stream of the binary data of the different layer

**Errors**:

- doesNotExist: If the artefact does not exist
- invalidArgument: If one of the input parameters is empty or has a wrong format
- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalOCIArtefact

**Description**: Deletes a local OCI artefact. 

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the local OCI artefact

**Outputs**:

**Errors**:

- existingReference: If the local OCI artefact is still referenced in some *Component Descriptor*
- invalidArgument: If one of the input parameters is empty or has a wrong format
- repositoryError: If some error occurred in the *Component Repository*

### getOciEndpointForLocalOciArtefact

**Description**: If the *Component Repository* supports the[Open Container Initiative Distribution Specification]
(https://github.com/opencontainers/distribution-spec/blob/main/spec.md) and provides a spec conforming endpoint
this endpoint could be retrieved by this function.

The v2 http endpoint to e.g. the manifest are as follows: 

``` registryEndPoint/v2/name/manifests/digestOfManifest ```

or 

``` registryEndPoint/v2/name/manifests/tag ```

The returned `name` needs not to be the same as specified in the parameter *ociArtefactName*, but it is recommended 
that *ociArtefactName* is at least a suffix of `<name>`.

**Inputs**:

- String ociArtefactName: name of the OCI artefact
- String tag: tag of the OCI artefact

**Outputs**:

- String registryEndPoint: endpoint of the OCI registry
- String name: name of the OCI artefact
- String digestOfManifest: digest of manifest

**Errors**:

- notSupported: If this kind of OCI access to the local OCI artefact is not supported by the *Component Repository*
- doesNotExist: If the local OCI artefact does not exist
- invalidArgument: If one of the input parameters is empty or has a wrong format
- repositoryError: If some error occurred in the *Component Repository*

**Example:**

Assume the method is called ass follows

```
    getOciEndpointForLocalOciArtefact("example.com/test", "0.1.1")
```

and returns the following data

```
    registryEndPoint="test.example-registry.com
    name="examplefolder/example.com/test"
    digestOfManifest="sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
```

then the manifest could be fetched using one of the following URLs

```
    test.example-registry.com/v2/examplefolder/example.com/test/manifests/sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623
```

or 

```
    test.example-registry.com/v2/examplefolder/example.com/test/manifests/0.1.1
```






