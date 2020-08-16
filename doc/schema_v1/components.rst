Resolvable Components and Resolvable Component References
=========================================================


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


Resolvable Component References (v1)
------------------------------------

Resolvable Component References consist exactly of the attributes `name`, `version`, where:

- `name` is the name of a Gardener Component
- `version` is version of a Gardener Component version

*Example*

.. code-block:: yaml

   name: 'component-name'
   version: '1.0.0'


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


Dependency Types
----------------

Any known type that is not a `Resolvable Component` is a `Dependency Type`. Dependencies
are used to describe technical artifacts the declaring resolvable component depends
upon. Dependencies are always defined by a `Resolvable Component`

Dependency Type Schema (v1)
---------------------------

Each `Dependency Types` *MUST* define the following attributes; `Dependencies` *MAY*
define additional components, according to their schema. Dependencies *MUST* only be
defined below the type-specific attribute below the `dependencies` attribute of a
`Resolvable Component`.

- name: str
- version: str, relaxed semver (see above)

*Example*

.. code-block:: yaml

   name: 'dependency-name'
   version: '1.0.0'


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

*Example (v1)*

.. code-block:: yaml

  dependencies:
    container_images:
      - name: 'example-image'
        version: '1.2.3'
        image_reference: 'eu.gcr.io/some-project/some-image:1.2.3'

*Example (v2)*

.. code-block:: yaml

  dependencies:
    - name: 'example-image'
      version: '1.2.3'
      type: 'ociImage'
      image_reference: 'eu.gcr.io/some-project/some-image:1.2.3'

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

*Example (v1)*

.. code-block:: yaml

  dependencies:
    web:
      - name: 'example-web-dependency'
        version: '1.2.3'
        url: 'https://example.org/some-file'


Generic Dependency Type
~~~~~~~~~~~~~~~~~~~~~~~

*since: `v1`*

An informal dependency intended for human interpretation.

+-----------------+----------------------------------------+
| name            | arbitrary name (ASCII recommened)      |
+-----------------+----------------------------------------+
| version         | relaxed semver                         |
+-----------------+----------------------------------------+

*Example (v1)*

.. code-block:: yaml

  dependencies:
    generic:
      - name: 'example-generic-dependency'
        version: '1.2.3'
