Gardener Component Descriptor Normalization
===========================================

Every resource can contain an optional digest attribut that specifies the digest of the referenced resource.
The digest is a object that *MUST* contain the `algorithm` and the hash as `value`.

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
        algorithm: 'sha256'
        value: 'abcdef'


Normalization
--------------

The normalized component descriptor has the following properties:
- Attributes *MUST* be transformed into list of single attributes (key-value pairs) and ordered alphabetically
- The component descriptor *MUST* be JSON encoded
- Resources contain their identity fields (name, version and extra identity) as well as their content digest.
- Sources are included only their identity fields
  - Sources without a dedicated artifact can be added to the normalized version by using the `extraIdentity` field.

.. code-block::
  - meta:
    - schemaVersion: 'v2'

  - component:
      - name: 'example.com/my-component'
      - version: 'v1.1.1'

      - sources:
        - - name: my-source
          - version: v0.0.0
          - extraIdentity: []

      - resources:
        - - name: my-artifact
          - version: v0.0.0
          - digest:
            - algorithm: sha256
            - value: 'abcdef'
          - extraIdentity:
            - abc: def


Digest/Fingerprint
------------------

The normalized component descriptor is not directly stored in the component but rather it is hashed with a specific algorithm.
The resulting digest is added to the component descriptor and persisted in the component registry.

Both signatures and the digest are appended to the component descriptor's `signatures` field.
.. code-block::
  meta:
    schemaVersion: 'v2'

  component:
    name: 'example.com/my-component'
    version: 'v1.1.1'

    resources:
    - name: my-artifact
      version: v0.0.0
      digest: 'sha256:abcdef'

  signatures:
  - normalizationType: v1
    digest:
      algo: sha256
      data: 'abcdef'
    signature:
      algo: RSA/DSA
      data: 'abcdef'