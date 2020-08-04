import os

own_dir = os.path.abspath(os.path.dirname(__file__))
repo_root = os.path.abspath(os.path.join(own_dir, os.pardir))
indep_dir = os.path.join(repo_root, 'language-independent')
test_res_dir = os.path.join(indep_dir, 'test-resources')
