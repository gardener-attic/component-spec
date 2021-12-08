Gardener Component Descriptor Normalization
===========================================

Every resource can contain an optional digest attribut that specifies the digest of the referenced resource.
The digest is an object that *MUST* contain the `hashAlgorithm` and the hash as `value`. Additionally, it *MUST* contain the normalisationAlgorithm.

.. code-block::
  meta:
    schemaVersion: 'v2'

  component:
    name: 'example.com/my-component'
    version: 'v1.1.1'

    resources:
    - name: my-artifact
      version: v0.0.0
      digest:
        hashAlgorithm: sha256
        normalisationAlgorithm: manifest-digest-v1
        value: 'abcdef'


Normalization
--------------

The normalized component descriptor has the following properties:
- Attributes *MUST* be transformed into list of single attributes (key-value pairs) and ordered alphabetically
- The component descriptor *MUST* be JSON encoded
- Component References are included with their identity fields (name, version and extra identity) and their content digest.
  The content digest contains hashAlgorithm, normalisationAlgorithm and value.
- Resources contain their identity fields (name, version and extra identity) as well as their content digest.
- Sources are included only their identity fields
  - Sources without a dedicated artifact can be added to the normalized version by using the `extraIdentity` field.

The following listing shows the list of single attributes in the normalised version.
.. code-block::
  - meta:
    - schemaVersion: 'v2'

  - component:
      - name: 'example.com/my-component'
      - version: 'v1.1.1'
      - digest:
        - hashAlgorithm: sha256
        - normalisationAlgorithm: sorted-alphabetically-v1
        - value: 'abcdef'

      - sources:
        - - name: my-source
          - version: v0.0.0
          - extraIdentity: []

      - resources:
        - - name: my-artifact
          - version: v0.0.0
          - digest:
            - hashAlgorithm: sha256
            - normalisationAlgorithm: manifest-digest-v1
            - value: 'abcdef'
          - extraIdentity:
            - abc: def


Digest/Fingerprint
------------------

The normalized component descriptor is not directly stored in the component but rather it is hashed with a specific algorithm.
The resulting digest is added to the component descriptor and persisted in the component registry.

Signature and digest are appended to the component descriptor's `signatures` field.
Multiple signatures can be addressed using the name field.
The normalisationAlgorithm defines the normalisation of the component-descriptor.

.. code-block::
  meta:
    schemaVersion: 'v2'

  component:
    name: 'example.com/my-component'
    version: 'v1.1.1'
    componentReferences:
    - name: 'example.com/my-component'
      version: 'v1.1.1'
      digest:
        algorithm: sha256
        normalisationAlgorithm: sorted-alphabetically-v1
        value: abcdef

    resources:
    - name: my-artifact
      version: v0.0.0
      extraIdentity:
        -key1: value1
      digest:
        algorithm: sha256
        normalisationAlgorithm: manifest-digest-v1
        value: abcdef

    sources:
    - name: source1
      version: v0.0.1
      extraIdentity:
      - key1: value1

  signatures:
  - name: signatureName
    digest:
      algorithm: sha256
      normalisationAlgorithm: sorted-alphabetically-v1
      value: abcdef
    signature:
      algorithm: RSA
      value: abcdef