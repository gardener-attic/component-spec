## Topics for Meeting with Uwe and Christian

Punkte für Component-Spec-Meeting
- Which resource types exists for cd.resource.access.type?
    - lokale (welche lokalen gibt es, wirklich nur einen)
    - was muss man für nicht lokale spezifizieren?
        - typspezifisches Format
        - welche Typen nehmen wir gleich auf
        - erweiterbar
- welche typen von sourcen gibt es git/http, welche sollen in die spec
- warum haben source entries 2 Type-Felder?
- Warum 3 Definitionen für Resourcen? https://gardener.github.io/component-spec/component-descriptor-v2.html#tab-pane_component_resources_items_anyOf_i2
- Was sollen diese Sätze?
    - In v2, for backwards-compatibility-reasons, the (mandatory) version attribute is implicitly added to the set of additional identity attributes iff it
      is not explicitly declared as such, and the uniqueness constraint would otherwise be violated.
      - https://gardener.github.io/component-spec/concepts.html
- Repository Context in CD
  - unglücklich, dass es kein Top-Level-Typfeld und der Rest in untergeordneter Struktur
  - seltsam, das ein Object auch seinen Speicherort beinhaltet, oder ist das nur für Transport-Historie?
  - aktuelles Schema erlaubt nur OCI-Ref. Wir haben aber auch Filesystem etc. und eine Standart muss hier erweiterbar sein!

## Todo
- Currently the access type for a local blob is localOciBlob, which sounds wrong. localBlob would be better.
- change technical artefacts to software artefacts 
- describe that name and version of a component descriptor are unique in the context of a component repository
- check json schema for component descriptor

## To check
---
Each of the artefacts declared in a `Component Descriptor` has an `identity`, an `access`
description and a `type`. `Artefact identities` can be used to reference an artefact in
the context of the declaring component descriptor. The `Artefact type` defines how the
artefact is to be interpreted. The `access Description` defines from where the artefact
can be retrieved.

When replicating component descriptors and their artefacts, the artefacts' `access` descriptions
may be changed. However, `artefact types` and `artefact identities` always remain unchanged.
---
Context Repositories
Component Descriptors are published to Context Repositories. Context repositories MUST always contain the transitive 
closure of all Component Descriptors, referenced by contained Comonent Descriptor versions.

Adhering to the aforementioned requirement of closure, Component Descriptors MAY be transported to other Context 
Repositories. When doing so, the history of Context Repositories MUST be retained, and appended. 
As part of such a transporting procedure, both the components’ sources and resources MAY be transferred to different 
repositories. The new locations MUST in this case be reflected in the new Component Descriptor.
---
If built from the declaring component’s sources, their versions MUST match the component’s version. Whether or not a 
resource is built from the referencing component is expressed through the relation attribute.

---

## Probably not needed

---
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
---
The CNUDIE Component model intends to solve the problem of addressing,
identifying, and accessing artefacts for software components, relative to an
arbitrary component repository. By that, it also enables the transport of
software components between component repositories.
---