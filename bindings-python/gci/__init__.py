import os

own_dir = os.path.abspath(os.path.dirname(__file__))

JSON_SCHEMA_DIR_NAME = 'jsonschema'
JSON_SCHEMA_FILE_NAME = 'component-descriptor-v2-schema.yaml'


def path_to_json_schema() -> str:
    if os.path.isfile(filepath := os.path.join(own_dir, JSON_SCHEMA_DIR_NAME, JSON_SCHEMA_FILE_NAME)):
        return filepath
    # fallback for local development
    elif os.path.isfile(
        filepath := os.path.join(own_dir, os.pardir, os.pardir, 'language-independent', JSON_SCHEMA_FILE_NAME)
    ):
        return filepath
    else:
        raise RuntimeError("Unable to find json schema file.")
