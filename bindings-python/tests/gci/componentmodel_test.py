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
