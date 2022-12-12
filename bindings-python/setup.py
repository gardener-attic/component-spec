import setuptools
import os

own_dir = os.path.abspath(os.path.dirname(__file__))


def version():
    if (version := os.environ.get('EFFECTIVE_VERSION')):
        return version

    print('warning: env var EFFECTIVE_VERSION not set - falling back to VERSION file')
    with open(os.path.join(own_dir, os.pardir, 'VERSION')) as f:
        return f.read().strip()


def requirements():
    with open(os.path.join(own_dir, 'requirements.txt')) as f:
        for line in f.readlines():
            line = line.strip()
            if not line or line.startswith('#'):
                continue
            yield line


setuptools.setup(
    name='gardener-component-model',
    version=version(),
    description='Gardener Component Model',
    python_requires='>=3.9.*',
    packages=setuptools.find_packages(),
    package_data={'gci': ['component-descriptor.json-schema.yaml']},
    install_requires=list(requirements()),
    entry_points={},
)
