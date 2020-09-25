#!/usr/bin/env python3

# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import os

import jsonschema
import yaml

own_dir = os.path.abspath(os.path.dirname(__file__))
repo_root = os.path.abspath(os.path.join(own_dir, os.path.pardir))
schema_file = os.path.join(repo_root, 'component-descriptor-v2-schema.yaml')


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--schema-file', default=schema_file)
    parser.add_argument('-f', '--component-descriptor', required=True)

    return parser.parse_args()


def main():
    parsed = parse_args()

    schema_file = parsed.schema_file
    component_descriptor = parsed.component_descriptor

    print(f'{schema_file=}')

    with open(schema_file) as f:
        schema_dict = yaml.safe_load(f)

    with open(component_descriptor) as f:
        comp_dict = yaml.safe_load(f)


    jsonschema.validate(
        instance=comp_dict,
        schema=schema_dict,
    )

    print('schema validation succeeded')


if __name__ == '__main__':
    main()
