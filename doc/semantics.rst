Gardener Component Descriptor Semantics
=======================================

`Component Descriptors` adhere to a schema defined in the JSON-schema contained in the
`component-spec <https://github.com/gardener/component-spec>`_ repository.

In this section, the semantics of those elements are defined.


Component
---------

`Components` are semantical entities with a certain semantical focus, as part
of a software product. They are typically built from `sources`, which are being
conciously developed.

Components have an identity that is defined through their name.
Component versions have an identity that is defined through component identity and version.

`Sources` are typically maintained in source code managment systems, and are
transformed into `Resources` (for example by a build), which are used at
installation or runtime for the product.

`Resources` are "local" if they are derived from a `source` declared by the same component.
`Resources` are "external" if they are not derived form a `source` declared by the same component.

For each of those entities, sets of versions exist. Over time, new versions are created
by changing sources.

Each component version defines a certain version vector or sources, from which a version vector
of resources can be derived.

Component Descriptor
~~~~~~~~~~~~~~~~~~~~

A `Component Descriptor version` is determined from a component version. It describes
the full version vectors of comprised source versions, and corresponding resource versions.

It *MAY* declare dependencies towards other component versions.

In addition, it may contain additional metadata.


Component Repository
~~~~~~~~~~~~~~~~~~~~

`Component Repositories` store `Component Descriptors versions` and allow referencing them through
the corresponding component version identities.

`Component Descriptor versions` stored in a component repository *SHOULD* be immutable.

`Component Repositories` *MAY* also be used to store sources or resources.


Component Repository History
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Components are typically being developed in a common context. Such a context consists of
organisational aspects (for example code ownership), and related infrastructure. Examples for
related infrastructure contain source code repositories, build infrastructure, and resource
repositories. An example for a source code repository being github.com, and an example for a
resource repository being an OCI Registry, such as eu.gcr.io for OCI images.

`Component Repositories` span such contexts.

Each component *MUST* have at least one such component repository, which is its
initial repository.

Additional component repositories *MAY* be defined, for example to model delivery
scenarios.

Component versions can be "transported" between component repositories by copying
the corresponding component descriptors, retaining their identities. This *MUST*
be done for the full transitive closure of referenced component versions.

In addition, referenced sources or resources *MAY* also be copied. If done,
the then-changed resource or source access information in the new (copied)
component descriptor versions *MUST* be adjusted.

In order to retain the transport history, the original repository history *MUST* be
retained. The target component repository must be appended to the list of component
repositories.

The last component repository is the "current" component repository.

Identities and Accessors
------------------------

Component versions are unambiguously identified by 2-tuples of `name` and `version` in the
context of a `Component Repository`. Their identity never changes across
any replication between context repositories.

Component descriptor versions reference `Source Versions`, `Resources` and other `Components`.

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
- names *MUST* start with a lowercase character

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
