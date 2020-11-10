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

Component versions are unambiguously identified by two-tuples of `name` and `version` in the
context of a `Component Descriptor Context Repository`. Their identity never changes across
any replication between context repositories.

Components declare dependencies towards `Sources`, `Resources` and other `Components`.

Any of those dependency declarations have an umambiguous identity in the context of their
declaring component version. Their identities always consist of their formal type
(`Source`, `Resource`, `Component Reference`) and a `name`.

Independent from the formal type, the following restrictions for valid names
are defined. All `name` attributes *MUST* be printable ASCII-characters from
any combination of the following character classes:

- lower-cased alphanumeric (`[a-z0-9]`)
- special characters (`[-_+]`)
- any other characters are *NOT* acceptable
- names *SHOULD* consist of at least two, and less than `64` characters

Dependency declarations *MAY* declare additional identity attributes.

Additional identity attribute names *MUST* adhere to the restrictions defined for `name` values
(see above). Their values *MUST* be UTF-8-encoded strings. It is strongly recommended to apply
the same restrictions, as defined for valid `names` (e.g. considering serialisation to
file-systems). As they are, however, component-specific, their concrete values are not assumed
to bear any concrete semantics, and are thus to be handled as opaque octet-sequences from any
generic tooling working on the component-model. Thus, any valid UTF-8 character is allowed.

Considering the tuple of `formal type`, `name`, and the optional additional identity attributes,
all dependency declarations *MUST* be unique for any given component version starting with the next
component descriptor version to be released after `v2`.

In `v2`, for backwards-compatibility-reasons, the (mandatory) `version` attribute is implicitly
added to the set of additional identity attributes iff it is not explicitly declared as such, and
the uniqueness constraint would otherwise be violated.

They also declare a format type (for example OCI Container Image).

Sources and resources also *MUST* declare through their `access` attribute a means to
access the underlying artifacts. This is done by declaring an access type (for example an OCI
Image Registry), which defines the protocol through which access is done. Depending on the
access type, additional attributes are required (e.g. an OCI Image Reference).

Component references do not need to declare an `access`, as the access method is defined by
the component descriptor specifiation.

In addition to the aforementioned mandatory attributes, any dependency declaration, and the
component itself may declare optional `labels`.


Component Sources
-----------------

Components are expected to have at least one source (code) representation. For each component
version, the corresponding source snapshot from which the component was built *MUST* be
specified.

Name Attribute
~~~~~~~~~~~~~~



Component References
--------------------

Component versions may declare dependencies towards other component versions. Dependencies are
always resolved in the same repository context. At the time of publishing a `Component Descriptor`,
`Component Descriptors` for each referenced component version *MUST* already be present in the
context repository.

A component descriptor registry *MUST* reject component descriptors with references to absent
component versions.

Resources
---------

resources are technical artifacts of a certain format (or `type`).

If built from the declaring component's sources, their `versions` *MUST* match the component's
version. Whether or not a resource is built from the referencing component is expressed through
the `relation` attribute.
