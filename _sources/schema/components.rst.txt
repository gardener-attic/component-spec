Resolvable Components and Resolvable Component References
=========================================================

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


Resolvable Component Schema (v1)
--------------------------------

Each `Resolvable Component` *MUST* define the following attributes; `Resolvable Components` *MAY*
define additional components, according to their schema.

- name: str, defined by Component Type
- version: str, relaxed semver (see above)
- dependencies: dict with dependencies (components, container_images, web, generic)

*Example*

.. code-block:: yaml

   name: 'component-name'
   version: '1.0.0'
   dependencies:
    components: []
    container_images: []
    web: []
    generic: []


Resolvable Component Schema (v2)
--------------------------------

Each `Resolvable Component` *MUST* define the following attributes; `Resolvable Components` *MAY*
define additional components, according to their schema.

- name: str, defined by Component Type
- version: str, relaxed semver (see above)
- type: str, one of `gardenerComponent`, `ociComponent`
- dependencies: list of dependencies (references to resolvable components or dependency components)

*Example*

.. code-block:: yaml

   name: 'component-name'
   version: '1.0.0'
   type: 'a-valid-component-type' # `gardenerComponent`|`ociComponent`
   dependencies: []


Resolvable Component References (v1)
------------------------------------

Resolvable Component References consist exactly of the attributes `name`, `version`, where:

- `name` is the name of a Gardener Component
- `version` is version of a Gardener Component version

*Example*

.. code-block:: yaml

   name: 'component-name'
   version: '1.0.0'

Resolvable Component References (v2)
------------------------------------

Resolvable Component References consist exactly of the attributes `name`, `version`, `type`, where:

- `name` is the name of a Gardener Component
- `version` is version of a Gardener Component version
- `type` is the referenced component's type (gardenerComponent or ociComponent)

*Example*

.. code-block:: yaml

   name: 'component-name'
   version: '1.0.0'
   type: 'gardenerComponent'


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

Schema (v1)
...........

Gardener Components do not define additional attributes.

*Example*

.. code-block:: yaml

   name: 'github.com/gardener/gardener'
   version: 'v1.7.2'
   dependencies:
    components:
      - name: 'github.com/gardener/etcd-druid'
        version: 'v0.3.0'
    container_images:
      - name: 'apiserver'
        version: 'v1.7.2'
        image_reference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'
    web: []
    generic: []

Schema (v2)
...........

Gardener Components do not define additional attributes.

*Example*

.. code-block:: yaml

   name: 'github.com/gardener/gardener'
   version: 'v1.7.2'
   type: 'gardenerComponent'
   dependencies:
    - name: 'github.com/gardener/etcd-druid'
      version: 'v0.3.0'
      type: 'gardenerComponent'
    - name: 'apiserver'
      version: 'v1.7.2'
      type: 'ociImage'
      image_reference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'

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
