Building Components
===================

The `CNUDIE Component Descriptor model` takes no assumption as to how component versions are created.
Typically, they are generated during a build process, as this is where their artifacts are generated.
About those build processes, the CNUDIE Component Descriptor model also does not take any
assumptions.

To make generating Component Descriptors more convenient, there is a re-usable toolset that can
be used (`<https://github.com/gardener/component-cli>`_).

The following example shows how this can be done based on a `make`-based build.


Example
-------

Prerequisites
~~~~~~~~~~~~~

- toolset must be available
- OCI Registry with write-credentials to store Component-Descriptors
- also: expose credentials (.docker/config.json)


Integrating Component Descriptor creation into the build
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

- todo: reference example build (w/o CNUDIE-CD)
- next step: run the build
- next step: add make-target + template-file to create CNUDIE-CD
- final step: run it
