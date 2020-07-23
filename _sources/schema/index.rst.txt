Component Descriptor specification
==================================

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
