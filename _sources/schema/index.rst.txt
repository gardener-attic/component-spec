Component Descriptor Schema specification
=========================================

.. toctree::
  :hidden:

  v1
  v2


Version-Independent specifications
----------------------------------

`Component Descriptors` are associative arrays of data, adhering to a versioned schema. While the
original implementation (now named `v1`) did not contain a version schema metadata, such an
attribute is introduced for `v2`, and will remain in place for future versions.

Metadata
~~~~~~~~

Each `Component Descriptor` of `v2` or later *MUST* contain a root-level attribute `meta`, which
*MUST* contain a nested attribute `schemaVersion`, whith the actual schema version as string value.

`Component Descriptors` of `v1` *MAY* contain such an attribute at root-level. If `meta` attribute
is absent, `v1` schema *MUST* be assumed.

*Example*

.. code-block:: yaml

   meta:
    schemaVersion: 'v2'


Component-Type
~~~~~~~~~~~~~~

Each resolvable component, and each dependency has a type. In `v1`, this type was implied by the
structure.  Since `v2`, the type is explcitly declared through an attribute `type`.

For each `Component Descriptor` schema version, there is a set of known types, with a well-defined
schema.

`Component Descriptors` *MUST* not contain any component or dependency entries with an unknown type.
Validation *MUST* fail in this case.

Custom Component-Type prefix
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To allow for custom extensions, the aforementioned validation rule does not apply to component types
with the reserved `x-` prefix. Component type names starting with the `x-` prefix *MUST* be retained
during deserialiation and serialisation operations, but *MAY* otherwise be ignored.


Uniqueness Constraints
~~~~~~~~~~~~~~~~~~~~~~

Each sequence of `Components` (either `Resolvable Components` or dependencies) *MUST NOT* contain
duplicate elements.

`Components` are duplicates, iff they reside in the same sequence, and share the same `name`,
`version`, and `type`.

Implementations for creation of `Component Descriptors` are _recommended_ to silently discard
duplicate entries if they are identical to an existing element, and *MUST* raise a
validation error if the duplicate entry definition differs.


Overwrite Semantics / Effective Component Descriptor
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In addition to declaring components and their dependencies, `Component Descriptors` *MAY* be
enriched with overwrite declarations (see version-specific schema definitions for details).

`Overwrites` are used, for example, to persist the result of a copying operation of a referenced
artifact to a new location, and thus provide a means to abstract the concrete location of a
dependency from the logical dependency.

Regardless of schema version, `overwrites` are always grouped in sets of individual overwrites,
representing one common overwrite operation.

Additional overwrite declarations *MAY* be appended to existing overwrite declarations. New
declarations *MUST* however always be appended as last sequence element. Existing overwrite
declarations *MUST NOT* be altered or omitted or re-ordered.

It is only permitted to overwrite commponent or dependency attributes that do not constitute to the
identification of the component or dependency. I.e. it is not permissible to overwrite `type`,
`name` or `version`.

To determine the `Effective Component Descriptor`, all overwrites *MUST* be applied in the order
they are specified to the defined components. All overwrites *MUST* only reference components
that were previously defined in the same `Component Descriptor`.

Consumers of dependencies declared in a `Component Descriptor` *SHOULD* always use the effective
dependencies.

.. vim: tw=100:
