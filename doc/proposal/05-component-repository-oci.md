# OCI Component Repository

This section describes an implementation of a Component Repository which stores component descriptors in an OCI image 
registry. We call this implementation *OCI Component Repository*.

An OCI component repository is defined by a *base URL*. This URL must refer to an OCI image registry or a path
into it. Every component descriptor in the repository is stored in the OCI image registry in an OCI artifact.
The name of the artifact is derived from the base URL, the component name, and component version:

```text
<base URL>/component-descriptors/<component name>:<component version>
```

If the resources of a component descriptor reference local OCI blobs, then these blobs are stored in the same OCI 
artifact as the component descriptor. The artifact consists of a manifest, a config, and an array of layers.
The component descriptor is the first layer. The local OCI blobs are stored in the other layers. 

![images/component-artifact.png](images/component-artifact.png)

The config of the artifact contains a reference to the component descriptor layer, as shown in the following example
of a config.json:

```json
{
    "componentDescriptorLayer": {
        "mediaType": "application/vnd.gardener.cloud.cnudie.component-descriptor.v2+yaml+tar",
        "digest": "sha256:6303ce21c1c8b5a9dfd9d6616cce976558d84542ac2eca342e3267f7205f2759",
        "size": 3584
    }
}
```