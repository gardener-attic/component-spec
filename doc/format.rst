Gardener Component Descriptor Format Specification
==================================================

`Gardener Component Descriptors` are associative arrays adhering to a versioned schema. As they are
intended to be a language-independent data format, also (de-)serialisation, publishing, discovery and
retrieval needs to be defined.

Data Format, Serialisation and Deserialisation
----------------------------------------------

In serialised form, `Component Descriptors` *MUST* be UTF-8-encoded. Either `YAML <https://yaml.org>`_, or
`JSON <https://json.org>`_ may be used. If `YAML` is used as serialisation format, only the subset of
features defined by `JSON` must be used, thus allowing conversion to a `JSON` representation.

`YAML` is recommended as preferred serialisation format.

`YAML` permits the usage of comments, and allows different formatting options. None of those are
by contract part of a `Gardener Component Descriptor`, thus implementations may arbitrarily choose
to retain or not retain comments or formatting options.

The order of attributes is insignificant, and *MUST NOT* be relied upon.

The order of elements in sequences *MAY* be significant and *MUST* be retained in cases where it
is significant.


Publishing and Discovery, Context Repositories
----------------------------------------------

Each `Component Version` is unambiguously identified by a two-tuple of:

- component name
- component version

Component Names
~~~~~~~~~~~~~~~

Component Names reside in a global namespace. To avoid name conflicts `Component Names` *MUST*
start with a valid domain name (as specified by RFC-1034) with an optional URL path
suffix (as specified by RFC-1738).

If no URL path suffix is specified, the domain *MUST* be possessed by the component proprietor.
If a URL path suffix is specified, the namespace started by the concatenation of domain and
URL path suffix *MUST* be possessed by the component proprietor.

The component name *SHOULD* reference a location where the component's resources (typically source
code, and/or documentation) are hosted.

An example, and recommended practise is using GitHub repository names for components on GitHub.

For example `github.com/gardener/gardener`.

Component Versions
~~~~~~~~~~~~~~~~~~

Component Versions refer to specific snapshots of a component. A common scenario being the release
of a component, thus offering it for consumption.

Component Versions *MUST* adhere to a loosened variant of `Semver 2.0.0 <https://semver.org>`_.

Different to strict semver 2.0.0, component versions *MAY*:

- have an optional `v` prefix
- omit the third component (patch-level); if omitted, path-level is implied to equal `0`

Context Repositories
~~~~~~~~~~~~~~~~~~~~

`Component Descriptors` are published to `Context Repositories`. Context repositories *MUST*
always contain the transitive closure of all `Component Descriptors`, referenced by contained
`Comonent Descriptor` versions.

Adhering to the aforementioned requirement of closure, `Component Descriptors` *MAY* be transported
to other `Context Repositories`. When doing so, the history of `Context Repositories` *MUST* be
retained, and appended.

As part of such a transporting procedure, both the components' sources and resources *MAY* be
transferred to different repositories. The new locations *MUST* in this case be reflected in the
new `Component Descriptor`.
