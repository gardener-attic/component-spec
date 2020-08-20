import dataclasses
import io
import os
import tarfile
import typing
import yaml

import gci.componentmodel


component_descriptor_fname = 'component-descriptor.yaml'


def component_descriptor_to_tarfileobj(
    component_descriptor: typing.Union[dict, gci.componentmodel.ComponentDescriptor],
):
    if not isinstance(component_descriptor, dict):
        component_descriptor = dataclasses.asdict(component_descriptor)

    component_descriptor_buf = io.BytesIO(
        yaml.dump(component_descriptor).encode('utf-8')
    )
    component_descriptor_buf.seek(0, os.SEEK_END)
    component_descriptor_leng = component_descriptor_buf.tell()
    component_descriptor_buf.seek(0)

    tar_buf = io.BytesIO()

    tf = tarfile.open(mode='w', fileobj=tar_buf)

    tar_info = tarfile.TarInfo(name=component_descriptor_fname)
    tar_info.size = component_descriptor_leng

    tf.addfile(tarinfo=tar_info, fileobj=component_descriptor_buf)
    tf.fileobj.seek(0)

    return tf.fileobj


def component_descriptor_from_tarfileobj(
    fileobj: io.BytesIO,
):
    with tarfile.open(fileobj=fileobj, mode='r') as tf:
        component_descriptor_info = tf.getmember(component_descriptor_fname)
        raw_dict = yaml.safe_load(tf.extractfile(component_descriptor_info).read())
        print(raw_dict)
        print('xx')

        gci.componentmodel.ComponentDescriptor.from_dict(raw_dict)
