# Introduction

The definition, structure and management of software in larger enterprises often builds upon tools and processes, which
largely originate from former on-premise thinking and monolithic architectures. Development teams responsible for
solutions or services have built specific, often point-2-point integrations with CI/CD systems, compliance tools,
reporting dashboards or delivery processes in the past. Larger development teams might have even built their own
toolsets specifically for their products, including those needed for compliance handling and delivery automation.
These concepts, process integrations and resulting tools are often still in use today, even though everyone knows:
They don't fit into today's cloud world.

The result is a fragmented set of homegrown specific tools across products, solutions and services, affecting an enterprises' ability to deliver
software consistently and compliant to its own or customer operated target environments. These specific, overly complex and thus hard to understand CI/CD pipelines, and the inability to instantly
provide a holistic aggregated view of currently running technical artifacts for each and every production environment
(including both cloud and on-premise), result in the overall management of software at scale becoming tedious, error-prone and
ineffective.

## Why is this a huge problem?

Most prominently, with the general unalignment of how software is defined and managed throughout the whole company,
it is not possible without additional overhead (like setting up even more processes and tools on top) to manage
the complete lifecycle of all solutions, services or individual deployment artifacts running in any
given landscape. Even worse, when trying to set up new landscapes, it becomes a nightmare to successfully orchestrate, deploy and configure the needed software components in the new environments.

As long as individual development teams within a company continue to use their own tools and processes to manage the
lifecycle of the software they are responsible for, this unsatisfying (and finally TCD and TCO affecting) situation can
not improve and will only get worse over time.

## How can this improve?
The major problem at hand here is the absence of one aligned software component model, consistently used across the
enterprise, to manage compliant software components and their technical artifacts. Such
a model would help not only with streamlined deployments to public and private cloud environments, but also in various
other areas of lifecycle management like compliance processes and reporting. This software component model must describe all technical artifacts of a software product, and establish an ID, which could then be used across all lifecycle management tasks.

First and foremost, the model has to be technology-agnostic at its heart, so that not only modern containerized cloud,
but also legacy software is supported, out-of-the-box. It simply has to be acknowledged that companies are not able to
just drop everything that has been used in the past and solely use shiny new cloud native workloads. This fact makes it
crucial to establish a common component model, which is able to handle both cloud native and legacy software, for which
it needs to be fully agnostic about the technology used.

Additionally, the model needs to be easily extensible. No one is able to
predict the future, apart from the fact that things will always change, especially in the area of IT. Being able to
adapt to future trends, without constantly affecting the processes and tools responsible for the core of the lifecycle
management of software, is a must.

**Todo: Describe why existing component models could not be used, why is our model better**
**Todo: Add image**

## Scope

Operating software installations/products, both for cloud and on-premises, covers many aspects:

- How, when and where are the technical artefacts created?
- How are technical artefacts stored and accessed?
- Which technical artefacts are to be deployed?
- How is the configuration managed?
- How and when are compliant checks, scanning etc. executed?
- When are technical artefacts deployed?
- Where and how are those artefacts deployed?
- Which other software installations are required and how are they deployed and accessed?
- etc.

The overall problem domain has a complexity that make it challenging to be solved as a whole.
However, the problem domain can be divided into two disjoint phases:

- production of technical artefacts
- deployment and lifecycle management of technical artefacts

The produced artefacts must be stored somewhere such that they can be accessed and collected for the deployment.
The OCM defines a standard to describe which technical artefacts belong to a software installation and how to
access them. This provides a clear interface between the production and the deployment/lifecycle management phase.

Though the following application areas are out of scope for OCM, it provides a common interface for
compliance checks, security scanning, code signing, transport, deployment or other lifecycle-management aspects.
If software installations are described using OCM, e.g. a scanning tool could use this to collect all technical
artefacts it needs to check. If the technical resources of different software installations are described with different
formalisms, such tools must provide interfaces and implementations for all if them and data exchange becomes a nightmare.

The problem becomes even harder if a software installation is build of different parts/components, each described with
another formalism. OCM allows a uniform definition of such dependencies such that one consistent description of 
a software installation is available.

The OCM does not make any assumptions about the (**Todo: Is this really the case?**)

- kinds of technical artefacts (e.g. docker images, helm chart, binaries etc., git sources)
- technology how to store and access technical artefacts (e.g. as OCI artefacts in an OCI registry)

OCM is and technology-agnostic specifications and allows implementations to define exactly those technical aspects
as an extension of the basic model.

**ToDo: Reference to our OCI based realization**

## Motivation Example

Usually complex software products are divided into logical units, which are called **components** in this specification.
For example, a software product might consist of three components, a frontend, a backend and some monitoring stack.
Of course, the software product itself could be seen as a component comprising the other three components.

As a result of the development phase, **component versions** are created, e.g. when you make a release for the component.

A component version consists of a set of technical artefacts, e.g. docker images, helm charts, binaries,
configuration data etc. Such artefacts are called **resources** in this specification. 

Resources are usually build from something, e.g. code in a git repo, which we call **sources** in this specification.

The OCM introduces a so called **Component Descriptor** for every component version, to describe the resources and sources 
belonging to a particular component version and how these could be accessed.

For the three components of our example software product, one *Component Descriptor* exists for every component version,
e.g. three *Component Descriptor* for the three versions of the frontend, six for the six versions of the backend etc.

Not all component version combinations of frontend, backend and monitoring result are compatible and build a product version.
In order to define reasonable version combinations for our software product we could use another feature of
the *Component Descriptor*, which allows the definition of dependencies to other component versions. 

For our example we could introduce a component for the overall product. A particular version of this product component 
is again described by a *Component Descriptor*, which contain dependencies to particular *Component Descriptors* for the
frontend, backend and monitoring.

This is only an example how to describe a particular product version with OCM as a component with one 
*Component Descriptor* with dependencies to other *Component Descriptors*, which itself could have dependencies and so on.
You are not restricted to this approach, i.e. you could still just maintain a list of component version combinations which
build a valid product release. But OCM provides you a simple approach to specify what belongs to a product version.
Starting with the *Component Descriptor* for a product version and following the component dependencies, you could
collect all artefacts, belonging to this product version.

Every *Component Descriptor* has a name and a version. 

*Component Descriptors* are stored in *Component Repositories*. A *Component Repository* must provide a possibility
to access every stored *Component Descriptor* only by its name and version. Another requirement to a
*Component Repository* is, that uploading a *Component Descriptors* fails if it has a dependency to another
*Component Descriptor* which is not contained in the *Component Repository*.

In a *Component Descriptor* dependencies to other *Component Descriptors* are expressed by their name and version. 
Also, the references to resources and sources in a *Component Descriptor* have (logical) names
(and some extra information explained later) as their unique identifier.

The concept of referencing other *Component Descriptors* only by their name and version allows the transport of
the *Component Descriptors* from one *Component Repository* to another without the need to change component
references in a *Component Descriptor* because the names and versions remain stable. Resources and sources
referenced by the *Component Descriptors* could either remain untouched or also transported to a new location. In the
latter case, the references in the transported *Component Descriptors* have to be updated to their new locations.

**Todo: Perhaps some small example image to make this more clear?**
