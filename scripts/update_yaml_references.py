"""
Script to update path references in YAML files after reorganizing the codebase.
"""

import os
import re
import sys
from pathlib import Path
import glob


def update_references_in_file(file_path, replacements):
    """Update path references in a YAML file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    for old_path, new_path in replacements.items():
        content = content.replace(old_path, new_path)
    
    if content != original_content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    
    return False


def find_yaml_files(directory, patterns):
    """Find all YAML files in a directory recursively matching the patterns."""
    yaml_files = []
    for pattern in patterns:
        for file_path in glob.glob(os.path.join(directory, pattern), recursive=True):
            if os.path.isfile(file_path) and (file_path.endswith('.yaml') or file_path.endswith('.yml')):
                yaml_files.append(file_path)
    return yaml_files


def main():
    """Main function."""
    repo_dir = sys.argv[1] if len(sys.argv) > 1 else '/home/ubuntu/repos/Fine-Tune'
    
    if not os.path.exists(repo_dir):
        print(f"Error: Repository directory {repo_dir} does not exist.")
        sys.exit(1)
    
    replacements = {
        "src/ml_infrastructure/kubeflow/manifests/": "k8s/kubeflow/manifests/",
        "src/ml_infrastructure/mlflow/config/": "k8s/mlflow/config/",
        "src/ml_infrastructure/kserve/manifests/": "k8s/kserve/manifests/",
        "src/ml_infrastructure/minio/manifests/": "k8s/minio/manifests/",
        "src/ml_infrastructure/minio/config/": "k8s/minio/config/",
        "src/ml_infrastructure/feast/config/": "k8s/feast/config/",
        "src/ml_infrastructure/seldon/manifests/": "k8s/seldon/manifests/",
        "src/ml_infrastructure/h2o/config/": "k8s/h2o/",
        "src/ml_infrastructure/jupyter/config/": "k8s/jupyter/",
        "argocd/": "k8s/argocd/",
        "monitoring/": "k8s/monitoring/"
    }
    
    yaml_files = find_yaml_files(repo_dir, ["**/*.yaml", "**/*.yml"])
    updated_files = 0
    
    for file_path in yaml_files:
        if update_references_in_file(file_path, replacements):
            print(f"Updated references in {file_path}")
            updated_files += 1
    
    print(f"\nSummary: Updated references in {updated_files} out of {len(yaml_files)} YAML files.")


if __name__ == '__main__':
    main()
