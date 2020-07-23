Resolvable Components
=====================

`Resolvable Components` in the context of the Gardener Project are parts from which Gardner
Landscapes can be materialised. They adhere to a common contract that allows for them to be
managed in a common and automated way.


Component Versions
------------------

Each `Resolvable Component` can be considered as an umbrella for a set of `Component Versions`.
Each `Resolvable Component Version` is a snapshot of said version, accompanied with a
`Component Descriptor`.

`Component Versions` are represented as (utf-encoded) strings adhering to a relaxed variant of
`semver <https://semver.org>`_.

Different to strict semver:

- an optional `v` prefix *MAY* be prepended
- the patch-component *MAY* be ommitted (defaults to `.0` in this case)

Component Types
---------------

Each `Resolvable Component` has a `Component Type`, which defines the contract through which
their `Component Versions` are published, discovered, and retrieved.


Gardener Component Type
~~~~~~~~~~~~~~~~~~~~~~~

`Gardener Components` are resolvable components that reside in a `GitHub <https://github.com>`_
repository. Each GitHub release is a component version and *MUST* contain a `Component Descriptor`
as GitHub release asset named `component_descriptor.yaml`.


OCI Component Type
~~~~~~~~~~~~~~~~~~

TBD
