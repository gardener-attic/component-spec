# Transport Format

We will define a transport format for component. 
As a subtask we first must define a transport format for different types of resources that can be referenced by 
components.

## Transport Format of an OCI Artifact

We describe a file system structure that represents a single OCI artifact.
A tar file of this file structure will be called the transport archive of the OCI artifact.

> Clarify that the containing directory is not part of the tar file.

The file system structure is a directory containing:

- an `artifact-descriptor.json` file
- and a `blobs` directory.

The `blobs` directory contains the manifest, config and the layer files of the OCI artifact as a flat file list.
In case of a multi arch artifact, the `blobs` directory can contain several (index) manifests and config files.
Every file has a filename according to its
[digest](https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests), 
where the algorithm separator character is replaced by a dot ("."). 

```text
artifact-archive
├── artifact-descriptor.json
└── blobs
    ├── sha.123... (manifest.json)
    ├── sha.234... (config.json)
    ├── sha.345... (layer)
    └── sha.456... (layer)
```

The `artifact-descriptor.json` contains the name of the artifact, and a mapping which associates with
the digest of a manifest a list of tags.

Example:

```json
{
  "name": "example-artifact",
  "digest-tag-mapping": [
    {
      "digest": "sha:123...", 
      "tags": ["v0.1.0"]
    }
  ]
}
```


## Transport Format for Other Resource Types

Whenever a new resource type is supported, a corresponding transport format must be defined.


## Transport Format of a Component

We describe a file system structure that represents a single component.
A tar file of this file structure will be called the transport archive of the component.

In a first step, all references to external resources are converted to resources of type `localOciBlob`.
This means two things: the resource in the component descriptor must be adjusted, and a local blob must be added.
We use the transport archive of the resource as local blob. 

We have already defined a representation of components as OCI artifacts. 
Therefore, we can use for components the same transport format as for OCI artifacts.

> Note that the transformation of external resources increases the number of layers.
> Hence the manifest of the original component (in its OCI representation) and the manifest in the transport format
> are different.

```text
artifact-archive
├── artifact-descriptor.json
└── blobs
    ├── sha.123... (manifest.json)
    ├── sha.234... (config.json)
    ├── sha.345... (component descriptor)
    ├── sha.456... (local blob / transport archive of a previously external resource)
    └── sha.567... (local blob / transport archive of a previously external resource)
```

The component version appears in the archive-descriptor.json as a tag associated to the digest of the component 
descriptor:

```json
{
  "name": "COMPONENT_NAME",
  "digest-tag-mapping": [
    {
      "digest": "sha:345...", 
      "tags": ["COMPONENT_VERSION"]
    }
  ]
}
```
