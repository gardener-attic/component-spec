# Local Blobs

A central problem, which should be resolved by OCM, is to provide a uniform and well-defined possibility to transport 
component versions from a source to a target location/landscape. Often, transports include additional 
intermediate locations/landscapes. The situation becomes even more complicated when particular locations are fenced, 
i.e. have no access to some external systems. Furthermore, often not only the *Component Descriptors* but also the 
referenced sources and resources must be transported. 

If some intermediate location does not provide a particular store for all of your artifact types (OCI 
images, helm charts...) and also has no access to the referenced artifacts, you need some mapping of these artifact 
types to upload them into the provided stores. Such a mapping has to be defined and implemented for all artifact type 
and store combinations. To reduce this overhead, *Component Repositories* MAY provide a possibility to store blobs in 
a type agnostic manner. Then you only need to define one blob format for every artifact type, which could be uploaded 
to every *Component Repository* providing such a functionality.

Another motivation for storing blobs in a *Component Repository* is, that often there are some additional
configuration data, you just want to store together with the *Component Descriptor* in an easy and uniform manner. 

Therefore, *Component Repositories* MAY provide the possibility to store technical artifacts together with the 
*Component Descriptors* in the *Component Repository* itself as so-called *local blobs*. This also allows to pack all 
component versions with their technical artifacts in a *Component Repository* as a completely self-contained package, a 
typical requirement if you need to deliver your product into a fenced landscape. 

As a short example, assume some component needs additional configuration data stored in some YAML file. If 
in some landscape of your transport chain there is only an OCI registry available to store content, then you need to 
define a format how to store such a YAML file in the OCI registry. With *local blobs* you could just upload the file into
the *Component Repository*. 

## Functions for Local Blobs

### UploadLocalBlob

If a *Component Repository* provides support for *local blobs* it MUST implement a method for uploading *local blobs*
as specified in this chapter.

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
      type: localBlob
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
      type: localBlob
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

If a *Component Repository* provides support for *local blobs* it MUST implement a method for fetching *local blobs*
as specified in this chapter.

**Description**: Fetches the binary data of a local blob. *localAccessInfo* is the *Component Repository* specific
access information you got when you uploaded the local blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String localAccessInfo: Access information of the local blob

**Outputs**:

- BinaryStream data: Binary stream containing the local blob data.

**Errors**:

- doesNotExist: If the local blob does not exist
- invalidArgument: If one of the input parameters is empty or invalid
- repositoryError: If some error occurred in the *Component Repository*

### ListLocalBlobs

If a *Component Repository* provides support for *local blobs* it MUST implement a method for listing *local blobs*
as specified in this chapter.

**Description**: Provides an iterator over all triples *componentName/componentVersion/localAccessInfo* of all
uploaded blobs. 

**Inputs**:

**Outputs**:

- Iterator over string triple: Triples of *componentName/componentVersion/localAccessInfo*

**Errors**:

- repositoryError: If some error occurred in the *Component Repository*

### DeleteLocalBlob

If a *Component Repository* provides support for *local blobs* it SHOULD implement a method for uploading *local blobs*
as specified in this chapter.

**Description**: Deletes a local blob. *localAccessInfo* is the *Component Repository* specific
information you got when you uploaded the local blob.

An error occurs if there is still an existing reference to the local blob.

**Inputs**:

- String name: Name of the *Component Descriptor*
- String version: Version of the *Component Descriptor*
- String localAccessInfo: Access information of the local blob

**Outputs**:

**Errors**:

- doesNotExist: If the local blob does not exist
- existingReference: If the local blob is still referenced
- invalidArgument: If one of the input parameters is empty
- repositoryError: If some error occurred in the *Component Repository*