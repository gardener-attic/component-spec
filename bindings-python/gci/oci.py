import dataclasses
import io
import tarfile
import typing
import yaml

import gci.componentmodel


def component_descriptor_to_tarfileobj(
    component_descriptor: typing.Union[dict, gci.componentmodel.ComponentDescriptor],
):
    if not isinstance(component_descriptor, dict):
        component_descriptor = dataclasses.asdict(component_descriptor)

    component_descriptor_buf = io.BytesIO(
        yaml.dump(component_descriptor).encode('utf-8')
    )

    tar_buf = io.BytesIO()

    tf = tarfile.open(mode='w', fileobj=tar_buf)

    tar_info = tarfile.TarInfo(name='component-descriptor.yaml')
    tar_info.size = component_descriptor_buf.tell()
    component_descriptor_buf.seek(0)

    tf.addfile(tarinfo=tar_info, fileobj=component_descriptor_buf)
    tf.fileobj.seek(0)

    return tf.fileobj
