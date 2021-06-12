Component Descriptor Registries
===============================

Component Descriptor Registries are used to store `Component Descriptors`. Using a
`Repository Context` configuration, and a `Component Name` and `Component Version`, it is possible
to unambiguously reference a `Component Descriptor` for the given component version.

Component Descriptor Registries *MUST* implement a WORM-semantics (write-once-read-many).

There can potentially be many different technical implementations. However, for the time being,
there is only one, using an OCI Image Registry.

Component Descriptor Registries backed by OCI Image Registry
------------------------------------------------------------

Component Descriptor Registries backed by an OCI Image Registry *MUST* specify a `baseUrl`. Said
URL *MUST* refer to an OCI Image Registry (for example `eu.gcr.io/gardener-project`).

The base URL *MAY* refer to sub-path on a OCI Image Registry.

The thus-referenced `Component Descriptor` is stored as an OCI Artifact. The
Image Manifest *MUST* at least blobs (one config-blob, and one layer-blob).

The config-blob *MUST* be a valid JSON-document and *MUST* reference the layer-blob containing
the serialised `Component Descriptor` by defining a toplevel-attribute `componentDescriptorLayer`,
which *MUST* contain the same attributes and values as the same blob-reference in the
OCI-Manifest.

The blob containing the `Component Descriptor` must be JSON- or YAML-serialised. It must be
`utf-8`-encoded. It *MUST* either contain:

- a POSIX.1-2001-compliant pax archive (commonly also known as "tarball"),
  containing as first entry a regular file named `component-descriptor.yaml`, containing
  the `Component Descriptor`
- the octet-sequence of the serialised `Component Descriptor`

Component Name Mapping
~~~~~~~~~~~~~~~~~~~~~~

Any OCI Image Registry used as a Component Descriptor registry *MUST* define an unambiguous
mapping of `Component Name`, `Component Version`-tuples to OCI Image References.

The following Component Name mappings are known:

- urlPath
- sha256-digest

Component Name Mapping `urlPath`
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Component Descriptor OCI Names are constructed like so:

`<baseUrl>/componentDescriptors/<componentName>:<versionTag>`

Where:

- baseUrl is the already-discussed registry URL
- componentName is a valid `Component Name` (i.e. a domain name with an optional URL path suffix)
- versionTag is a valid `Component Version` (i.e. relaxed semver version)


Component Name Mapping `sha256-digest`
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Component Descriptor OCI Names are constructed like so:

`<baseUrl>/<componentNameDigest>:<versionTag>`

Where:

- baseUrl is the already-discussed registry URL
- ComponentName is the lowercased hexadecimal sha256-digest of the utf-8-encoded `Component Name`
- versionTag is a valid `Component Version` (i.e. relaxed semver version)
