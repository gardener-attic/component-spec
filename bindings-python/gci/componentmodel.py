import dataclasses
import enum
import functools
import io
import logging
import typing
import urllib.parse

import dacite
import jsonschema
import yaml

from . import path_to_json_schema

dc = dataclasses.dataclass

logger = logging.getLogger(__name__)


class ValidationMode(enum.Enum):
    FAIL = 'fail'
    WARN = 'warn'
    NONE = 'none'


class SchemaVersion(enum.Enum):
    V1 = 'v1'
    V2 = 'v2'


class AccessType(enum.Enum):
    OCI_REGISTRY = 'ociRegistry'
    GITHUB = 'github'
    HTTP = 'http'
    LOCAL_OCI_BLOB = 'localOciBlob'
    LOCAL_FILESYSTEM_BLOB = 'localFilesystemBlob'
    NONE = 'None'  # the resource is only declared informally (e.g. generic)

class SourceType(enum.Enum):
    GIT = 'git'


class ResourceType(enum.Enum):
    OCI_IMAGE = 'ociImage'
    GENERIC = 'generic'
    BLUEPRINT = 'blueprint'
    HELM = 'helm'

class ResourceRelation(enum.Enum):
    LOCAL = 'local'
    EXTERNAL = 'external'


@dc(frozen=True)
class ResourceAccess:
    type: typing.Union[AccessType, str]


@dc(frozen=True)
class OciAccess(ResourceAccess):
    imageReference: str


@dc(frozen=True)
class GithubAccess(ResourceAccess):
    repoUrl: str
    ref: str
    commit: typing.Optional[str] = None

    def __post_init__(self):
        parsed = self._normalise_and_parse_url()
        if not len(parsed.path[1:].split('/')):
            raise ValueError(f'{self.repoUrl=} must have exactly two path components')

    @functools.lru_cache
    def _normalise_and_parse_url(self):
        parsed = urllib.parse.urlparse(self.repoUrl)
        if not parsed.scheme:
            # prepend dummy-schema to properly parse hostname and path (and rm it again later)
            parsed = urllib.parse.urlparse('dummy://' + self.repoUrl)
            parsed = urllib.parse.urlunparse((
                '',
                parsed.netloc,
                parsed.path,
                '',
                '',
                '',
            ))
            parsed = urllib.parse.urlparse(parsed)

        return parsed

    def repository_name(self):
        return self._normalise_and_parse_url().path[1:].split('/')[1]

    def org_name(self):
        return self._normalise_and_parse_url().path[1:].split('/')[0]

    def hostname(self):
        return self._normalise_and_parse_url().hostname


@dc(frozen=True)
class HttpAccess(ResourceAccess):
    url: str

@dc(frozen=True)
class LocalOCIBlobAccess(ResourceAccess):
    digest: str

@dc(frozen=True)
class LocalFilesystemBlobAccess(ResourceAccess):
    filename: str
    mediaType: typing.Optional[str] = None


@dc(frozen=True)
class Label:
    name: str
    value: typing.Union[str, int, float, bool, dict]


_no_default = object()


class FindLabelMixin:
    def find_label(
        self,
        name: str,
        default=_no_default,
        raise_if_absent: bool = False,
    ):
        for label in self.labels:
            if label.name == name:
                return label
        else:
            if default is _no_default and raise_if_absent:
                raise ValueError(f'no such label: {name=}')
            if default is _no_default:
                return None
            return default


class Provider(enum.Enum):
    '''
    internal: from repositoryContext-owner
    external: from 3rd-party (not repositoryContext-owner)
    '''
    INTERNAL = 'internal'
    EXTERNAL = 'external'


@dc(frozen=True)
class Metadata:
    schemaVersion: SchemaVersion = SchemaVersion.V2


class ArtifactIdentity:
    def __init__(self, name, **kwargs):
        self.name = name
        kwargs['name'] = name
        # ensure stable order to ensure stable sort order
        self._id_attrs = tuple(sorted(kwargs.items(), key=lambda i: i[0]))

    def __str__(self):
        return '-'.join((a[1] for a in self._id_attrs))

    def __len__(self):
        return len(self._id_attrs)

    def __eq__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs == other._id_attrs

    def __hash__(self):
        return hash((type(self), self._id_attrs))

    def __lt__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs.__lt__(other._id_attrs)

    def __le__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs.__le__(other._id_attrs)

    def __ne__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs.__ne__(other._id_attrs)

    def __gt__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs.__gt__(other._id_attrs)

    def __ge__(self, other):
        if not type(self) == type(other):
            return False
        return self._id_attrs.__ge__(other._id_attrs)


class ComponentReferenceIdentity(ArtifactIdentity):
    pass


class ResourceIdentity(ArtifactIdentity):
    pass


class SourceIdentity(ArtifactIdentity):
    pass


@dc(frozen=True)
class ComponentIdentity:
    name: str
    version: str


class Artifact:
    '''
    base class for ComponentReference, Resource, Source
    '''
    def identity(self, peers: typing.Sequence['Artifact']):
        '''
        returns the identity-object for this artifact (component-ref, resource, or source).

        Note that, in component-descriptor-v2, the `version` attribute is implicitly added iff
        there would otherwise be a conflict, iff this artifact only uses its `name` as
        identity-attr (which is the default).

        In future versions of component-descriptor, this behaviour will be discontinued. It will
        instead be regarded as an error if the IDs of a given sequence of artifacts (declared by
        one component-descriptor) are not all pairwise different.
        '''
        own_type = type(self)
        for p in peers:
            if not type(p) == own_type:
                raise ValueError(f'all peers must be of same type {type(self)=} {type(p)=}')

        if own_type is ComponentReference:
            IdCtor = ComponentReferenceIdentity
        elif own_type is Resource:
            IdCtor = ResourceIdentity
        elif own_type is ComponentSource:
            IdCtor = SourceIdentity
        else:
            raise NotImplementedError(own_type)

        identity = IdCtor(
            name=self.name,
            **(self.extraIdentity or {})
        )

        if not peers:
            return identity

        if len(identity) > 1:  # special-case-handling not required if there are additional-id-attrs
            return identity

        # check whether there are collissions
        for peer in peers:
            if peer is self:
                continue
            if peer.identity(peers=()) == identity:
                # there is at least one collision (id est: another artifact w/ same name)
                return ArtifactIdentity(
                    name=self.name,
                    version=self.version,
                )
        # there were no collisions
        return identity


@dc(frozen=True)
class ComponentReference(Artifact, FindLabelMixin):
    name: str
    componentName: str
    version: str
    extraIdentity: typing.Dict[str, str] = dataclasses.field(default_factory=dict)
    labels: typing.List[Label] = dataclasses.field(default_factory=tuple)


@dc(frozen=True)
class SourceReference(FindLabelMixin):
    identitySelector: typing.Dict[str, str]
    labels: typing.List[Label] = dataclasses.field(default_factory=tuple)


@dc(frozen=True)
class Resource(Artifact, FindLabelMixin):
    name: str
    version: str
    type: typing.Union[ResourceType, str]
    access: typing.Union[
        OciAccess,
        GithubAccess,
        HttpAccess,
        LocalOCIBlobAccess,
        LocalFilesystemBlobAccess,
        ResourceAccess,
        None,
    ]
    extraIdentity: typing.Dict[str, str] = dataclasses.field(default_factory=dict)
    relation: ResourceRelation = ResourceRelation.LOCAL
    labels: typing.List[Label] = dataclasses.field(default_factory=tuple)
    srcRefs: typing.List[SourceReference] = dataclasses.field(default_factory=tuple)


@dc
class RepositoryContext:
    baseUrl: str
    type: AccessType = AccessType.OCI_REGISTRY


@dc
class ComponentSource(Artifact, FindLabelMixin):
    name: str
    access: GithubAccess
    version: typing.Optional[str] = None  # introduce this backwards-compatible for now
    extraIdentity: typing.Dict[str, str] = dataclasses.field(default_factory=dict)
    type: SourceType = SourceType.GIT
    labels: typing.List[Label] = dataclasses.field(default_factory=list)


@dc
class Component(FindLabelMixin):
    name: str     # must be valid URL w/o schema
    version: str  # relaxed semver

    repositoryContexts: typing.List[RepositoryContext]
    provider: Provider

    sources: typing.List[ComponentSource]
    componentReferences: typing.List[ComponentReference]
    resources: typing.List[Resource]

    labels: typing.List[Label] = dataclasses.field(default_factory=list)

    def current_repository_ctx(self):
        if not self.repositoryContexts:
            return None
        return self.repositoryContexts[-1]

    def identity(self):
        return ComponentIdentity(name=self.name, version=self.version)


@functools.lru_cache
def _read_schema_file(schema_file_path: str):
    with open(schema_file_path) as f:
        return yaml.safe_load(f)


def enum_or_string(v, enum_type: enum.Enum):
  try:
    return enum_type(v)
  except ValueError:
    return str(v)


@dc
class ComponentDescriptor:
    meta: Metadata
    component: Component

    @staticmethod
    def validate(
        component_descriptor_dict: dict,
        validation_mode: ValidationMode,
        json_schema_file_path: str = None,
    ):
        if validation_mode is ValidationMode.NONE:
            return

        json_schema_file_path = json_schema_file_path or path_to_json_schema()
        schema_dict = _read_schema_file(json_schema_file_path)

        try:
            jsonschema.validate(
                instance=component_descriptor_dict,
                schema=schema_dict,
            )
        except jsonschema.ValidationError as e:
            if validation_mode is ValidationMode.WARN:
                logger.warn(f'Error when validating Component Descriptor: {e}')
            elif validation_mode is ValidationMode.FAIL:
                raise
            else:
                raise NotImplementedError(validation_mode)

    @staticmethod
    def from_dict(
        component_descriptor_dict: dict,
        validation_mode: ValidationMode = ValidationMode.NONE,
    ):
        component_descriptor = dacite.from_dict(
            data_class=ComponentDescriptor,
            data=component_descriptor_dict,
            config=dacite.Config(
                cast=[
                    AccessType,
                    Provider,
                    ResourceType,
                    SchemaVersion,
                    SourceType,
                    ResourceRelation,
                ],
                type_hooks={
                    typing.Union[AccessType, str]: functools.partial(enum_or_string, enum_type=AccessType),
                    typing.Union[ResourceType, str]: functools.partial(enum_or_string, enum_type=ResourceType),
                },
            )
        )
        if not validation_mode is ValidationMode.NONE:
            ComponentDescriptor.validate(
                component_descriptor_dict=component_descriptor_dict,
                validation_mode=validation_mode,
            )

        return component_descriptor

    def to_fobj(self, fileobj: io.BytesIO):
        raw_dict = dataclasses.asdict(self)
        yaml.dump(
            data=raw_dict,
            stream=fileobj,
            Dumper=EnumValueYamlDumper,
        )


class EnumValueYamlDumper(yaml.SafeDumper):
    '''
    a yaml.SafeDumper that will dump enum objects using their values
    '''
    def represent_data(self, data):
        if isinstance(data, enum.Enum):
            return self.represent_data(data.value)
        return super().represent_data(data)
