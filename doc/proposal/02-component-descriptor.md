# Component Descriptor Specification

*Component Descriptors* are the central concept of OCM. A *Component Descriptor* describes what belongs to a particular
version of a software component and how to access it. This includes:

- resources, i.e. technical artefacts like binaries, docker images, ...
- sources like code in github
- dependencies to other software component versions

## Component Descriptor Format Specification

A *Component Descriptor* is a [YAML](https://yaml.org/) or [JSON](https://www.json.org/json-en.html) document
according to this [schema](component-descriptor-v2-schema.yaml).

In serialised form, Component Descriptors MUST be UTF-8-encoded. Either YAML, or JSON may be used. If YAML is used
as serialisation format, only the subset of features defined by JSON must be used, thus allowing conversion to a
JSON representation.

YAML is recommended as preferred serialisation format.

YAML permits the usage of comments, and allows different formatting options. None of those are by contract part of a
*Component Descriptor*, thus implementations may arbitrarily choose to retain or not retain comments or formatting
options.

The order of attributes is insignificant, and MUST NOT be relied upon.

The order of elements in sequences MAY be significant and MUST be retained in cases where it is significant.

## Schema Version

A *Component Descriptor* document consists of two top level elements: *meta*, *component*

The *meta* element contains the schema version of the *Component Descriptor* specification. This document defines
schema version *v2*.

```
meta:
  - schemaVersion: "v2"
component:
  ...
```

## Name and Version

Every *Component Descriptor* has a name (*component.name*) and version (*component.version), also called component
name and component version. Name and version are the identifier for a *Component Descriptor* and the component version 
described by it.

```
meta:
  - schemaVersion: "v2"
component:
  name: ...
  version: ...
```

Component names reside in a global namespace. To avoid name conflicts component names MUST start with a valid domain
name (as specified by RFC-1034, RFC-1035) with an optional URL path suffix (as specified by RFC-1738).

If no URL path suffix is specified, the domain MUST be possessed by the component proprietor. If a URL path suffix is
specified, the namespace started by the concatenation of domain and URL path suffix MUST be possessed by the component
owner.

The component name SHOULD reference a location where the component’s resources (typically source code, and/or
documentation) are hosted. An example and recommended practise is using GitHub repository names for components on
GitHub like *github.com/gardener/gardener*.

Component versions refer to specific snapshots of a component. A common scenario being the release of a component.
Component versions MUST adhere to a loosened variant of [Semver 2.0.0](https://semver.org/).

Different to strict semver 2.0.0, component versions MAY:

- have an optional v prefix
- omit the third component (patch-level); if omitted, path-level is implied to equal 0

## Sources, Resources

Components versions are typically built from sources, maintained in source code management systems,
and transformed into resources (for example by a build), which are used at installation or runtime of the product.

Each *Component Descriptor* contains a field for references to the used sources and a field for references
to the used resources.

Example for sources:

```
...
component:
  sources:
  - name: example-sources-1
    version: v1.19.4
    type: git
    access:
      commit: e01326928b6f9825dba9fa530b8d4917f93194b0
      ref: refs/tags/v1.19.4
      repoUrl: github.com/gardener/example-sources-1
      type: github
      ...
```

Every source entry has a name and version field. Furthermore, it has a type field, specifying the nature of the
source code management system, and a type specific access section.

Example for resources:

```
...
component:
  resources:
    - name: external-monitoring
      version: v0.8.3
      relation: external
      type: ociImage
      access:
        imageReference: example.com/monitoring:v0.8.3
        type: ociRegistry
    
    - name: example-chart
      version: v1.19.4
      relation: local
      type: helm.io/chart
      access:
        digest: sha256:b5733194756a0a4a99a4b71c4328f1ccf01f866b5c3efcb4a025f02201ccf623
        type: localOciBlob
        mediaType: oci.image+tar
      ...
```

Every resource entry has a *name* and *version* field. The *type* of the resource specifies the content of the
referenced resource, i.e. if it is a helm chart, a json file etc.

Resources have a *relation* field with value “local” if they are derived from a source declared by the same component.
The value is “external” if they are not derived form a source declared by the same component, e.g. a docker image for
a monitoring stack developed in another component. 

If the *relation* field has the value “local”, the *version* field of the resource reference must be the same as the 
*version* field of the *Component Descriptor*.

The *access* field of a resource entry contains the information how to access the resource. Every *access* entry has a
required field *type* which specifies the access method.

Sources and resources declare through their access attribute a means to access the underlying artifacts.
This is done by declaring an access type (for example an OCI Image Registry), which defines the protocol through which
access is done. Depending on the access type, additional attributes are required (e.g. an OCI Image Reference).

OCM does not specify the format and semantics of particular access types for resources and sources. This can be done 
in extensions to this specification. An exception are the two access types  *localOciBlob* and *localOciArtefact* 
explained later.

## Component References

A component version might have dependencies to other component versions. The semantics of component version dependencies
is not defined here. Such dependencies could be expressed for different relationships between component versions,
e.g. deploy order, a component version comprise other component versions, a component version uses artefacts of
another component version etc.

To express dependencies to other component versions, a *Component Descriptor* has a field to specify dependencies to 
other *Component Descriptors*.

```
component:
  componentReferences:
  - name: name-1
    componentName: github.com/.../component-name-1
    version: v1.38.3
  - name: name-2
    componentName: .../component-name-2
    version: v0.11.4 
```

Every component reference has a *name* field. This name is only the identifier of the component reference within
the *Component Descriptor*. Furthermore, every entry contains the component name and version of the referenced
*Component Descriptor*.

As later elaborated in the context of *Component Repositories*, component references do not need to declare an access
attribute.

## Identifier for Sources, Resources and Component References

Every entry in the *sources*, *resources* and *componentReferences* fields has a *name* field. The following restrictions
for valid names are defined:

- lower-cased alphanumeric ([a-z0-9])
- special characters ([-_+])
- any other characters are NOT acceptable
- names SHOULD consist of at least two, and less than 64 characters
- names MUST start with a lowercase character

Every *source*, *resource* or *componentReference* needs a unique identifier in a *Component Descriptor*.
In particular situations the name is not sufficient, e.g. if docker images for different platform are included.
Therefore, every entry has an additional optional field *extraIdentity* to resolve this problem, i.e. every entry
must have a unique combination of *name*, *extraIdentity* and formal type (*source*, *resource* or *componentReference*)
within a *Component Descriptor*.

An *extraIdentity* is a map, of key value pairs whereby:

- The keys must adhere to the same restrictions defined for name values (see above)
- The values MUST be UTF-8-encoded strings.

Two *extraIdentities* are equal iff they have the same key value pairs whereby the order is not relevant.

Example for twi resource entries with the same name but different extra identities and therefore different identifier:

```
component:
  resources:
  - name: name-1
    extraIdentity:
      platform: "arm64"
      country: "us"
    ...
  - name: name-1
    extraIdentity:
      platform: "x86_64"
      country: "de"
    ...
```

## Labels for Sources, Resources and Component References

Every entry in the *sources*, *resources* and *componentReferences* fields, and the component itself may declare
optional labels. This allows to express application specific extensions.

## Repository Contexts

Every *Component Descriptor* has a field *repositoryContexts* containing an array of access information of 
*Component Descriptor Repositories*, i.e. a stores for *Component Descriptors* which are specified later. 

The array of access information describes the transport chain of a *Component Descriptor* through different
*Component Descriptor Repositories*, whereby the last entry describes the current *Component Descriptor Repository*,
in which the *Component Descriptor* is stored.

## Misc

Other fields of a *Component Descriptor* are:

- component.provider: provider of the component,e.g. a company, organization,...

- component.resources.srcRefs: References to resources have another field *srcRefs*. If the corresponding resource was build 
from "local" sources these could be listed here by providing their identifier within the *Component Descriptor*, i.e. 
their names and extraIdentities.