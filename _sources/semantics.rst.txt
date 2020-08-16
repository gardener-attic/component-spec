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

tbd

Component Sources
-----------------

tbd

Component References
--------------------

tbd

Local Resources
---------------

tbd

External Resources
------------------

tbd
