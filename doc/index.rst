.. Gardener Component-Descriptor Contract and Specification documentation master file, created by
   sphinx-quickstart on Thu Jul 23 08:36:05 2020.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Open Component Model
====================


This document describes the Open Component Model (OCM), its idea and reasoning, scope, contract
and some tools working on top of this model.

The component model provides a machine-readable way to describe versioned artefact sets that finally
build installable software packages, that can be transported among public and local repository contexts.
It is a sound technology-independent basis for cooperating tools to access artefacts in a
locaction-agnostic manner. This can be used to run tools and processes uniformly for the
same software in various, even fenced environments.

.. toctree::
  :maxdepth: 2
  :caption: Contents:

  motivation
  concepts
  format
  component_descriptor_registries
  oci
  building_components


Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
