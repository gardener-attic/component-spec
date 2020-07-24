Component Descriptor Schema specification
=========================================

.. toctree::
  :hidden:

  v1
  v2


Version-Independent specifications
----------------------------------

`Component Descriptors` are associative arrays of data, adhering to a versioned schema. While the original implementation
(now named `v1`) did not contain a version schema metadata, such an attribute is introduced for `v2`, and will remain in
place for future versions.

Metadata
~~~~~~~~

Each `Component Descriptor` of `v2` or later *MUST* contain a root-level attribute `meta`, which *MUST* contain a nested
attribute `schema_version`, whith the actual schema version as string value.

`Component Descriptors` of `v1` *MAY* contain such an attribute at root-level. If `meta` attribute is absent, `v1` schema
*MUST* be assumed.

*Example*

.. code-block:: yaml

   meta:
    schema_version: 'v2'


Component-Type
~~~~~~~~~~~~~~

Each resolvable component, and each dependency has a type. In `v1`, this type was implied by the structure.
Since `v2`, the type is explcitly declared through an attribute `type`.

For each `Component Descriptor` schema version, there is a set of known types, with a well-defined schema.

`Component Descriptors` *MUST* not contain any component or dependency entries with an unknown type.
Validation *MUST* fail in this case.

Custom Component-Type prefix
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To allow for custom extensions, the aforementioned validation rule does not apply to component types with the
reserved `x-` prefix. Component type names starting with the `x-` prefix *MUST* be retained during
deserialiation and serialisation operations, but *MAY* otherwise be ignored.


Uniqueness Constraints
~~~~~~~~~~~~~~~~~~~~~~

Each sequence of `Components` (either `Resolvable Components` or dependencies) *MUST NOT* contain
duplicate elements.

`Components` are duplicates, iff they reside in the same sequence, and share the same `name`,
`version`, and `type`.

Implementations for creation of `Component Descriptors` are _recommended_ to silently discard
duplicate entries if they are identical to an existing element, and *MUST* raise a
validation error if the duplicate entry definition differs.
