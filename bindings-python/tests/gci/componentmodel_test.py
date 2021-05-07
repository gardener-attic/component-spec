import dataclasses
import os
import typing
import unittest

import dacite
import yaml

import gci.componentmodel as cm
import paths

test_res_dir = paths.test_res_dir


def test_deserialisation():
    with open(os.path.join(test_res_dir, 'component_descriptor_v2.yaml')) as f:
        component_descriptor_dict = yaml.safe_load(f)
    component_descriptor = cm.ComponentDescriptor.from_dict(
        component_descriptor_dict=component_descriptor_dict,
    )

    assert component_descriptor.component.resources[0].type is cm.ResourceType.OCI_IMAGE
    assert isinstance(component_descriptor.component.resources[0].access, cm.OciAccess)
    assert component_descriptor.component.resources[0].access.type is cm.AccessType.OCI_REGISTRY


def test_deserialisation_of_custom_resources():
    with open(os.path.join(test_res_dir, 'component_descriptor_v2_custom.yaml')) as f:
        component_descriptor_dict = yaml.safe_load(f)

    component_descriptor = cm.ComponentDescriptor.from_dict(
        component_descriptor_dict=component_descriptor_dict,
    )
    assert isinstance(component_descriptor.component.resources[0].access, cm.LocalFilesystemBlobAccess)
    assert component_descriptor.component.resources[1].access is None
    assert isinstance(component_descriptor.component.resources[2].access, cm.ResourceAccess)


def test_github_access():
    gha = cm.GithubAccess(
        repoUrl='github.com/org/repo',
        ref='refs/heads/master',
        type=cm.AccessType.GITHUB,
    )

    assert gha.repository_name() == 'repo'
    assert gha.org_name() == 'org'
    assert gha.hostname() == 'github.com'


def test_component():
    component = cm.Component(
        name='component-name',
        version='1.2.3',
        repositoryContexts=[
            cm.RepositoryContext(baseUrl='old-ctx-url'),
            cm.RepositoryContext(baseUrl='current-ctx-url'),
        ],
        provider=None,
        sources=(),
        componentReferences=(),
        resources=(),
        labels=(),
    )

    assert component.current_repository_ctx().baseUrl == 'current-ctx-url'


def test_patch_label():
    lssd_label_name = 'cloud.gardener.cnudie/sdo/lssd'
    processing_rule_name = 'test-processing-rule'

    @dataclasses.dataclass
    class TestCase(unittest.TestCase):
        name: str
        input_labels: typing.List[cm.Label]
        label_to_patch: cm.Label
        raise_if_absent: bool
        expected_labels: typing.List[cm.Label]
        expected_err_msg: str

    testcases = [
        TestCase(
            name='appends label to empty input_labels list',
            input_labels=[],
            label_to_patch=cm.Label(
                name=lssd_label_name,
                value={
                    'processingRules': [
                        processing_rule_name,
                    ],
                },
            ),
            raise_if_absent=False,
            expected_labels=[
                cm.Label(
                    name=lssd_label_name,
                    value={
                        'processingRules': [
                            processing_rule_name,
                        ],
                    },
                ),
            ],
            expected_err_msg=''
        ),
        TestCase(
            name='throws exception if len(input_labels) == 0 and raise_if_absent == True',
            input_labels=[],
            label_to_patch=cm.Label(
                name=lssd_label_name,
                value={
                    'processingRules': [
                        processing_rule_name,
                    ],
                },
            ),
            expected_labels=None,
            raise_if_absent=True,
            expected_err_msg=f'no such label: name=\'{lssd_label_name}\'',
        ),
        TestCase(
            name='throws no exception if label exists and raise_if_absent == True',
            input_labels=[
                cm.Label(
                    name=lssd_label_name,
                    value={
                        'processingRules': [
                            'first-pipeline',
                        ],
                    },
                ),
            ],
            label_to_patch=cm.Label(
                name=lssd_label_name,
                value={
                    'processingRules': [
                        processing_rule_name,
                    ],
                },
            ),
            raise_if_absent=True,
            expected_labels=[
                cm.Label(
                    name=lssd_label_name,
                    value={
                        'processingRules': [
                            processing_rule_name,
                        ],
                    },
                ),
            ],
            expected_err_msg=''
        ),
        TestCase(
            name='overwrites preexisting label',
            input_labels=[
                cm.Label(
                    name='test-label',
                    value='test-val',
                ),
                cm.Label(
                    name=lssd_label_name,
                    value={
                        'processingRules': [
                            'first-pipeline',
                        ],
                        'otherOperations': 'test',
                    },
                ),
            ],
            label_to_patch=cm.Label(
                name=lssd_label_name,
                value={
                    'processingRules': [
                        processing_rule_name,
                    ],
                },
            ),
            raise_if_absent=False,
            expected_labels=[
                cm.Label(
                    name='test-label',
                    value='test-val',
                ),
                cm.Label(
                    name=lssd_label_name,
                    value={
                        'processingRules': [
                            processing_rule_name,
                        ],
                    },
                ),
            ],
            expected_err_msg='',
        ),
    ]

    for testcase in testcases:
        test_resource = cm.Resource(
            name='test-resource',
            version='v0.1.0',
            type=cm.ResourceType.OCI_IMAGE,
            access=cm.OciAccess(
                cm.AccessType.OCI_REGISTRY,
                imageReference='test-repo.com/test-resource:v0.1.0'
            ),
            labels=testcase.input_labels,
        )

        if testcase.expected_err_msg != '':
            with testcase.assertRaises(ValueError) as ctx:
                patched_resource = test_resource.set_label(
                    label=testcase.label_to_patch,
                    raise_if_absent=testcase.raise_if_absent,
                )
            assert testcase.expected_err_msg == str(ctx.exception)
        else:
            patched_resource = test_resource.set_label(
                label=testcase.label_to_patch,
                raise_if_absent=testcase.raise_if_absent,
            )
            testcase.assertListEqual(
                list1=patched_resource.labels,
                list2=testcase.expected_labels,
            )
