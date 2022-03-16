# Component Repository Specification

*Component Descriptors* are stored in *Component Repositories*. This chapter describes in a programing language/protocol
agnostic fashion the functions a *Component Repository* MUST implement to be OCM conform.
These functions are also specifying the interface whit which the user can access a *Component Repository* to perform 
operations on *Component Descriptors*.

Usually sources and resources are stored in some external storage, e.g. a docker image is stored in some OCI registry
and the *Component Descriptor* contains only the information how to access it. As an alternative *Component Repositories* 
also provide the possibility to store technical artifacts together with the *Component Descriptors* in the 
*Component Repository* itself as so-called *local blobs*. This allows to pack all component versions with their 
technical artifacts in a *Component Repository* as a completely self-contained package. This is a typical requirement 
if you need to deliver your product into a fenced landscape. This also allows storing e.g. configuration data together 
with your *Component Descriptor*.

**Todo: Image mit CDs, localblobs ...**

## Functions for Component Descriptors

### UploadComponentDescriptor

**Description**: Uploads a *Component Descriptor* to the *Component Repository*. If successful, the *Component Descriptor*
is accessible by its name and version (see GetComponentDescriptor). The name and version of a *Component Descriptor*
is the identifier of a *Component Descriptor*, therefore if there already exists a *Component Descriptor*
with the same name and version, the upload fails. 

A *Component Repository* MUST check if all referenced *Component Descriptors* and local blobs are already stored in 
the *Component Repository*.

If the identifier of entries in *resources*, *sources* or *componentReferences* are not unique, as described before,
an *invalidArgument* error MUST be returned.

If the identifier of *resources* and *sources* (name plus extraIdentity) are not unique, a *Component Repository* might 
automatically add the version field to the extraIdentity to resolve this problem. Of course, it MUST still fail, if 
uniqueness of the resource identifiers could not be achieved this way.

If the last entry in the "repositoryContext" field is not an entry for the current *Component Repository*, such an 
entry is automatically added.

**Inputs**:

- String componentDescriptor: JSON or YAML representation of the *Component Descriptor*

**Outputs**:

- Bool: true if upload was successful 

**Errors**:

- alreadyExists: If there already exists a *Component Descriptor* with the same name and version
- missingReference: If a referenced *Component Descriptor* or local blob does not exist in the *Component Repository*
- invalidArgument: If the parameter *componentDescriptor* is missing or is not conform to the
  [specified json schema](component-descriptor-v2-schema.yaml)
- repositoryError: If some error occurred in the *Component Repository*

### GetComponentDescriptor

**Description**: Returns the *Component Descriptor* as a YAML or JSON string according the 
[JSON schema](component-descriptor-v2-schema.yaml).

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*

**Outputs**:

- String: YAML or JSON string of the *Component Descriptor*

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

**Description**: Allows uploading binary data. The binary data belong to a particular *Component Descriptor* 
and can be referenced by the component descriptor in its *resources* or *sources* section. 
*Component Descriptors* are not allowed to reference local blobs of other *Component Descriptors* in their resources. 

When uploading a local blob, it is not REQUIRED that the corresponding *Component Descriptor* already exists. 
Usually local blobs are uploaded first because it is not allowed to upload a *Component Descriptor* if its local 
blobs not already exist.

The optional parameter *mediaType* provides information about the internal structure of the provided blob.

With the optional parameter *annotations* you could provide additional information about the blob. This information
could be used by the *Component Repository* itself or later if the local blob is stored again in some external
location, e.g. an OCI registry.   

*LocalAccessInfo* provides the information how to access the blob data with the method *GetLocalBlob* (see below). 

With the return value *globalAccessInfo*, the *Component Repository* could optionally provide an external reference to 
the resource, e.g. if the blob contains the data of an OCI image it could provide an external OCI image reference. 

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- BinaryStream data: Binary stream containing the local blob data.
- String mediaType: media-type of the uploaded data (optional)
- map(string,string) annotations: Additional information about the uploaded artifact (optional)

**Outputs**:

- String localAccessInfo: The information how to access the source or resource as a *local blob*.
- String globalAccessInfo (optional): The information how to access the source or resource via a global reference.  

**Errors**:

- invalidArgument: If one of the input parameters is empty or not valid
- repositoryError: If some error occurred in the *Component Repository*

**Example**: 
Assume you want to upload an OCI image to your *Component Repository* with the *UploadLocalBlob* function with media type
*application/vnd.oci.image.manifest.v1+json* and the *annotations* "name: test/monitoring", and get the *localAccessInfo*: 

```
"digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623"
```

Then the entry in the *Component Descriptor* might look as follows. It is up to you, if you add the annotations 
provided to the upload function and depends on the use case.

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
... 
```

The *Component Repository* could also provide some *globalAccessInfo* containing the location in an OCI registry:

```
imageReference: somePrefix/test/monitoring@sha:...
type: ociRegistry
```

An entry to this resource with this information in the *Component Descriptor* looks as the following:

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

### GetLocalBlob

**Description**: Fetches the binary data of a local blob. *blobIdentifier* is the *Component Repository* specific 
access information you got when you uploaded the local blob.

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

**Description**: Deletes a local blob. *blobIdentifier* is the *Component Repository* specific
information you got when you uploaded the local blob.  

An error occurs if there is still an existing reference to the local blob.


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
