import os

import dacite
import yaml

import gci.componentmodel as cm

own_dir = os.path.abspath(os.path.dirname(__file__))

with open(os.path.join(own_dir, 'component_descriptor_v2.yaml')) as f:
    component_descriptor_dict = yaml.safe_load(f)


def test_deserialisation():
    component_descriptor = dacite.from_dict(
        data_class=cm.ComponentDescriptor,
        data=component_descriptor_dict,
        config=dacite.Config(
            cast=[
                cm.SchemaVersion,
                cm.ComponentType,
            ]
        )
    )
