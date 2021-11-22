Motivation and Scope
====================

Operating Software installations both for Cloud and on-premises covers many aspects:

- how, when and where are the technical artefacts created
- how are technical artefacts stored and accessed
- which technical artefacts are to be deployed
- how is the configuration managed
- when are technical artefacts deployed
- where and how are those artefacts deployed
- which other software installations are required and how are they deployed and accessed

The overall problem domain has a complexity that make it challenging to be solved as a whole.
However, the problem domain can be divided into two disjoint subdomains:

- production of software artefacts
- deployment and lifecycle management of installations

The elements transitioning between those two subdomains are synchronisation
data ("when things should happen") and technical artefacts and their purposes.

By standardising those elements, it is possible to decouple those two aspects. This allows
to split the overall complexity into smaller parts. By this, the overall complexity can also be
reduced.

The CNUDIE Component Model defines a standardisation of the latter (technical
artefacts and their purposes). The major goal that is tried to be achieved by
this standardisation is to allow for components and their artefacts to be
umambiguously addressed and accessed in a location- and technology-agnostic
manner. This should work for both global and local or even private environments.


Core Concepts
-------------

:ref:`Components <concepts_components>` in CNUDIE are semantical entities with a
certain semantical focus, as part of a software product. They are typically
built from :ref:`Sources <concepts_sources`, which are being conciously developed.

They consist of a set of (component) versions. Each version is described by a
:ref:`Component Descriptor <concepts_component_descriptors>`, which describes
the included set of artefacts. `Component Descriptors` are stored in a
:ref:`Component Repository <_concepts_component_repositories` in a standardised
way. Within such a component repository, component descriptors are addressed by
their component name and version.

By separating the storage locations (component repository) from the component version
identity, in combination with a defined addressing scheme, it is possible to
access component descriptors in a component repository, in a location-independent manner.
This can be leveraged to replicate component descriptors between component repositories.

Each of the artefacts declared in a `Component Descriptor` has an `identity`, an `access`
description and a `type`. `Artefact identities` can be used to reference an artefact in
the context of the declaring component descriptor. The `Artefact type` defines how the
artefact is to be interpreted. The `access Description` defines from where the artefact
can be retrieved.

When replicating component descriptors and their artefacts, the artefacts' `access` descriptions
may be changed. However, `artefact types` and `artefact identities` always remain unchanged.

The CNUDIE Component Model can roughly be compared to a filesystem "living" in an arbitrary
storage (e.g. a blobstore, oci-registry, archive, ..).

Starting from an arbitrary component repository used as root repository it is
possible to access any artefact described by the component model as long as the
dedicated component version has been imported into this repository. The
dedicated component repository together with the component version identity
allows to access the component descriptor which then is used to determine the
access for dedicated artefacts given by their local identity. This way the
component respository acts as composed filesystem with fixed top level folders,
the components and their versions, followed by the artefact level, which then
may point into any another (kind of) repository. Sub folders can be described
by component version references described by component descriptors, which again
feature a local identity in the describing component descriptor.


Scope Definition
----------------

The CNUDIE Component model intends to solve the problem of addressing,
identifying, and accessing artefacts for software components, relative to an
arbitrary component repository. By that, it also enables the transport of
software components between component repositories.

Through the standardisation of structure and `access` to artefacts, it can serve as an
interface to any operations that need to interact with the content. This allows
for tools operating against this interface to be implemented in a re-usable, and
technology-agnostic manner (examples being transport, compliance and security
scanning, codesigning, ..).

Higherlevel functionality, such as deployment, or lifecycle-management aspects are
out of the scope the CNUDIE Component Model targets. Data for such aspects (for example
a deployment blueprint) can, however, be described by the CNUDIE Component Model
as dedicated typed artefacts. This is equivalent to e.g. adding a `Makefile`
into a filesystem. The filesystem does not "know" about the semantics of the
`Makefile` (including, e.g. declared dependencies towards other files).
