import dataclasses
import enum
import functools
import io
import typing
import urllib.parse

import dacite
import yaml

dc = dataclasses.dataclass


class SchemaVersion(enum.Enum):
    V1 = 'v1'
    V2 = 'v2'


class AccessType(enum.Enum):
    OCI_REGISTRY = 'ociRegistry'
    GITHUB = 'github'
    HTTP = 'http'
    NONE = 'None' # the resource is only declared informally (e.g. generic)


class SourceType(enum.Enum):
    GIT = 'git'


class ResourceType(enum.Enum):
    OCI_IMAGE = 'ociImage'
    GENERIC = 'generic'


class ResourceRelation(enum.Enum):
    LOCAL = 'local'
    EXTERNAL = 'external'


@dc(frozen=True)
class ResourceAccess:
    type: AccessType


@dc(frozen=True)
class OciAccess(ResourceAccess):
    imageReference: str


@dc(frozen=True)
class GithubAccess(ResourceAccess):
    repoUrl: str
    ref: str
    type: AccessType
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
class Label:
    name: str
    value: typing.Union[str, int, float, bool, dict]


_no_default = object()


class FindLabelMixin:
    def find_label(
        self,
        name: str,
        default=_no_default,

    ):
        for label in self.labels:
            if label.name == name:
                return label
        else:
            if default is _no_default:
                raise ValueError(f'no such label: {name=}')
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



@dc(frozen=True)
class ComponentReference(FindLabelMixin):
    name: str
    componentName: str
    version: str
    labels: typing.List[Label] = dataclasses.field(default_factory=list)


@dc(frozen=True)
class Resource(FindLabelMixin):
    name: str
    version: str
    type: ResourceType
    access: typing.Union[
        OciAccess,
        GithubAccess,
        HttpAccess,
        ResourceAccess,
        None
    ]
    relation: ResourceRelation = ResourceRelation.LOCAL
    labels: typing.List[Label] = dataclasses.field(default_factory=list)


@dc
class RepositoryContext:
    baseUrl: str
    type: AccessType = AccessType.OCI_REGISTRY


@dc
class ComponentSource(FindLabelMixin):
    name: str
    access: typing.Union[
        GithubAccess,
    ]
    type: SourceType = SourceType.GIT
    labels: typing.List[Label] = dataclasses.field(default_factory=list)


@dc
class Component(FindLabelMixin):
    name: str    # must be valid URL w/o schema
    version: str # relaxed semver

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


@dc
class ComponentDescriptor:
    meta: Metadata
    component: Component

    @staticmethod
    def from_dict(component_descriptor_dict: dict):
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
                ]
            )
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
