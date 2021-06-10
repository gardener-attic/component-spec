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

The thus-referenced `Component Descriptor` is stored as an OCI Container Image would be. The
Image Manifest *MUST* contain exactly one layer. The layer *MUST* be a POSIX.1-2001-compliant
pax archive (commonly also known as "tarball"), containing as first entry a regular file named
`component-descriptor.yaml`, whose `utf-8`-encoded content is a valid `Component Descriptor` (v2),
serialised as either YAML or JSON.

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
