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

The order of attributes is insignificant, and must not be relied upon.


Publishing and Discovery
------------------------

For any `Resolvable Component Version`, it *MUST* be possible to unambiguously discover and
retrieve a `Component Descriptor` using the three-tuple of:

- component name
- component version
- component type

See :doc:`component` for component-type specifics.
