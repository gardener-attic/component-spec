Gardener Component Descriptor Semantics
=======================================

`Component Descriptors` adhere to a schema defined in the JSON-schema contained in the
`component-spec <https://github.com/gardener/component-spec>`_ repository.

In this section, the semantics of those elements are defined.


Repository Contexts History
---------------------------

Sets of Components are typically being developed in a common context. Such a context consists of
organisational aspects (for example code ownership), and related infrastructure. Examples for
related infrastructure contain source code repositories, build infrastructure, and resource
repositories. An example for a source code repository being github.com, and an example for a
resource repository being an OCI Registry, such as eu.gcr.io for OCI images.

Such contexts are accompanied by `Component Descriptor Context Repositories`, which are used to
manage the component metadata (in the form of `Component Descriptors`).

Each component *MUST* have at least one such repository context, which is its initial context.

For delivery scenarios, additional repository contexts *MAY* be defined. Such a repository context
consists of the same set of technical repository types. In particular, it *MUST* contain a
separate `Component Descriptor` repository.

Component versions can be "transported" into such a context, by copying all referenced resources
into the corresponding resource repositories of the target context. In addition to copying
referenced resources, a new `Component Descriptor` *MUST* be created, containing references to
the new resource locations.

In order to retain the transport history, the original repository context history *MUST* be
retained, and be appended to with the new (target) repository context.


Provider
--------

Each component has a provider, which is declared relative to the original repository context.
In most cases, components are "internal", which means that the entity maintaining the component
is the same as the proprietor of the repository context.

Opposed to that, components maintained by a third party are declared through the value `external`.

Identities and Accessors
------------------------

Sources and Resources are unambiguously identified by two-tuples of name and version, in the
context of their declaring component version.

They also declare a format type (for example OCI Container Image).

Finally, sources and resources also *MUST* declare through their `access` attribute a means to
access the underlying artifacts. This is done by declaring an access type (for example an OCI
Image Registry), which defines the protocol through which access is done. Depending on the
access type, additional attributes are required (e.g. an OCI Image Reference).

Component Sources
-----------------

Components are expected to have at least one source (code) representation. For each component
version, the corresponding source snapshot from which the component was built *MUST* be
specified.

Component References
--------------------

Component versions may declare dependencies towards other component versions. Dependencies are
always resolved in the same repository context. At the time of publishing a `Component Descriptor`,
`Component Descriptors` for each referenced component version *MUST* already be present in the
context repository.

A component descriptor registry *MUST* reject component descriptors with references to absent
component versions.

Local Resources
---------------

Local resources are technical artifacts of a certain format (or type) that are built from the
declaring component's sources. Their versions *MUST* always match the version of their component.

External Resources
------------------

External resources differ from local resources in that they are _required_ by the declaring
component, but not built from it. Their versions *SHOULD* follow the versioning of their provided
(upstream) artifacts. External resources are a means to express dependencies towards the technical
artifacts of an external component that does not abide by the `Component Descriptor` contract.
