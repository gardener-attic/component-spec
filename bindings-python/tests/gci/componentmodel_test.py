# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

import os

import dacite
import yaml

import gci.componentmodel as cm
import paths

test_res_dir = paths.test_res_dir


with open(os.path.join(test_res_dir, 'component_descriptor_v2.yaml')) as f:
    component_descriptor_dict = yaml.safe_load(f)


def test_deserialisation():
    component_descriptor = cm.ComponentDescriptor.from_dict(
        component_descriptor_dict=component_descriptor_dict,
    )
