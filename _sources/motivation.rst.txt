Motivation and Scope
====================

Operating Software installations both for Cloud and on-premises covers many aspects:

- how, when and where are the technical artifacts created
- how are technical artifacts stored and accessed
- which technical artifacts are to be deployed
- how is the configuration managed
- when are technical artifacts deployed
- where and how are those artifacts deployed
- which other software installations are required and how are they deployed and accessed

The overall problem domain has a complexity that make it challenging to be solved as a whole.
However, the problem domain can be divided into two disjoint subdomains:

- production of software artifacts
- deployment and lifecycle management of installations

The elements transitioning between those two subdomains are synchronisation
data ("when things should happen") and technical artifacts and their purposes.

By standardising those elements, it is possible to decouple those two aspects. This allows
to split the overall complexity smaller parts. By this, the overall complexity can also be
reduced.

The CNUDIE Component Descriptor model defines a standardisation of the latter (technical
artifacts and their purposes).
