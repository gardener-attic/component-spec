import dataclasses
import enum
import typing

import dacite

dc = dataclasses.dataclass


class SchemaVersion(enum.Enum):
    V1 = 'v1'
    V2 = 'v2'


class ComponentType(enum.Enum):
    GARDENER_COMPONENT = 'gardenerComponent'
    OCI_COMPONENT = 'ociComponent'

    OCI_IMAGE = 'ociImage'
    WEB_DEPENDENCY = 'web'
    GENERIC = 'generic'


@dc(frozen=True)
class Metadata:
    schemaVersion: SchemaVersion


@dc(frozen=True)
class ComponentReference:
    name: str
    version: str
    type: ComponentType


@dc(frozen=True)
class ResolvableComponentReference(ComponentReference):
    pass


@dc(frozen=True)
class DependencyComponent(ComponentReference):
    pass


@dc(frozen=True)
class OciImage(DependencyComponent):
    imageReference: str


@dc(frozen=True)
class WebDependency(DependencyComponent):
    url: str


@dc(frozen=True)
class GenericDependency(DependencyComponent):
    pass


@dc
class ResolvableComponent:
    dependencies: typing.List[
        typing.Union[
            ComponentReference,
            DependencyComponent,
        ]
    ]


@dc
class GardenerComponent(ResolvableComponent):
    pass


@dc
class OciComponent(ResolvableComponent):
    pass


@dc
class Overwrite:
    componentReference: ResolvableComponentReference
    componentOverwrites: dict = dataclasses.field(
        default_factory=dict,
    ) # XXX define explicit overwrite-type(s)
    dependencyOverwrites: typing.List[
        typing.Union[
            ResolvableComponentReference,
            OciImage,
            WebDependency,
            GenericDependency,
        ]
    ] = dataclasses.field(default_factory=list)


@dc
class OverwriteDeclaration:
    declaringComponent: ResolvableComponentReference
    overwrites: typing.List[Overwrite]


@dc
class ComponentDescriptor:
    meta: Metadata
    components: typing.List[
        typing.Union[
            GardenerComponent,
            OciComponent,
        ]
    ]
    overwriteDeclarations: typing.List[OverwriteDeclaration]

    @staticmethod
    def from_dict(component_descriptor_dict: dict):
        component_descriptor = dacite.from_dict(
            data_class=ComponentDescriptor,
            data=component_descriptor_dict,
            config=dacite.Config(
                cast=[
                    SchemaVersion,
                    ComponentType,
                ]
            )
        )

        return component_descriptor
