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
their `Component Versions` are published, discovered, and retrieved. Each resolvable component
version may be declared as a dependency by other resolvable components.


Gardener Component Type
~~~~~~~~~~~~~~~~~~~~~~~

*since: `v1`*

`Gardener Components` are resolvable components that reside in a `GitHub <https://github.com>`_
repository. Each GitHub release is a component version and *MUST* contain a `Component Descriptor`
as GitHub release asset named `component_descriptor.yaml`.

+----------------------+-------------------------------------------+
| component name       | GitHub repository URL (without schema)    |
+----------------------+-------------------------------------------+
| component version    | GitHub release tag                        |
+----------------------+-------------------------------------------+
| component descriptor | release asset `component_descriptor.yaml` |
+----------------------+-------------------------------------------+


OCI Component Type
~~~~~~~~~~~~~~~~~~

*since: `v2`*

TBD

+----------------------+----------------------------------------+
| component name       | OCI reference without tag              |
+----------------------+----------------------------------------+
| component version    | OCI reference tag                      |
+----------------------+----------------------------------------+
| component descriptor | GitHub release tag                     |
+----------------------+----------------------------------------+


Dependency Types
----------------

Any known type that is not a `Resolvable Component` is a `Dependency Type`. Dependencies
are used to describe technical artifacts the declaring resolvable component depends
upon.


OCI Image Dependency Type
~~~~~~~~~~~~~~~~~~~~~~~~~

*since: `v1`*

An OCI Container Image published to an OCI Image registry.


+-----------------+----------------------------------------+
| name            | arbitrary name (ASCII recommened)      |
+-----------------+----------------------------------------+
| version         | relaxed semver                         |
+-----------------+----------------------------------------+
| image_reference | OCI image reference                    |
+-----------------+----------------------------------------+


Web Dependency Type
~~~~~~~~~~~~~~~~~~~

*since: `v1`*

A dependency retrievable via HTTP-GET

+-----------------+----------------------------------------+
| name            | arbitrary name (ASCII recommened)      |
+-----------------+----------------------------------------+
| version         | relaxed semver                         |
+-----------------+----------------------------------------+
| url             | a Unified Resource Locator             |
+-----------------+----------------------------------------+


Generic Dependency Type
~~~~~~~~~~~~~~~~~~~~~~~

*since: `v1`*

An informal dependency intended for human interpretation.

+-----------------+----------------------------------------+
| name            | arbitrary name (ASCII recommened)      |
+-----------------+----------------------------------------+
| version         | relaxed semver                         |
+-----------------+----------------------------------------+
