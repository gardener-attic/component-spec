import setuptools
import os

own_dir = os.path.abspath(os.path.dirname(__file__))


def requirements():
    with open(os.path.join(own_dir, 'requirements.txt')) as f:
        for line in f.readlines():
            line = line.strip()
            if not line or line.startswith('#'):
                continue
            yield line


setuptools.setup(
    name='gardener-component-model',
    version='0.0.0alpha4',
    description='Gardener Component Model',
    python_requires='>=3.8.*',
    packages=setuptools.find_packages(),
    install_requires=list(requirements()),
    entry_points={
    },
)
