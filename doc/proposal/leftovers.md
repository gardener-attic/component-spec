## Topics for Component-Spec-Meeting

- localOciBlob is a bad name, localBlob is much better

- Repository Context in CD
  - current schema only allows OCI-Ref, this needs to be extensible
  - Should OCI ref be part of this spec?
  
- CTF:
  - Should be included in OCM?
  
- What are the next step to specify?
  - particular types
  - Transport, CTF
  - Signing 
  - etc.

## Todo
- change technical artefacts to software artefacts 
- references in the text to the CD spec for particular entries

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