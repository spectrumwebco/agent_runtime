"""
Script to update import statements in Python files after reorganizing the codebase.
"""

import os
import re
import sys
from pathlib import Path


def update_imports_in_file(file_path):
    """Update import statements in a Python file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    replacements = {
        r'from\s+src\.ml_infrastructure\.api\.models\b': 'from src.models.api.ml_infrastructure_api_models',
        r'from\s+src\.ml_infrastructure\.data\.validation\.models\b': 'from src.models.data_validation.validation_models',
        r'from\s+src\.integrations\.github\.models\b': 'from src.models.integrations.github_models',
        r'from\s+src\.integrations\.gitee\.models\b': 'from src.models.integrations.gitee_models',
        r'from\s+src\.integrations\.issue_collector\.models\b': 'from src.models.integrations.issue_collector_models',
        r'import\s+src\.ml_infrastructure\.api\.models\b': 'import src.models.api.ml_infrastructure_api_models',
        r'import\s+src\.ml_infrastructure\.data\.validation\.models\b': 'import src.models.data_validation.validation_models',
        r'import\s+src\.integrations\.github\.models\b': 'import src.models.integrations.github_models',
        r'import\s+src\.integrations\.gitee\.models\b': 'import src.models.integrations.gitee_models',
        r'import\s+src\.integrations\.issue_collector\.models\b': 'import src.models.integrations.issue_collector_models'
    }
    
    for pattern, replacement in replacements.items():
        content = re.sub(pattern, replacement, content)
    
    if content != original_content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    
    return False


def find_python_files(directory):
    """Find all Python files in a directory recursively."""
    python_files = []
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.py'):
                python_files.append(os.path.join(root, file))
    return python_files


def main():
    """Main function."""
    repo_dir = sys.argv[1] if len(sys.argv) > 1 else '/home/ubuntu/repos/Fine-Tune'
    src_dir = os.path.join(repo_dir, 'src')
    
    if not os.path.exists(src_dir):
        print(f"Error: Source directory {src_dir} does not exist.")
        sys.exit(1)
    
    python_files = find_python_files(src_dir)
    updated_files = 0
    
    for file_path in python_files:
        if update_imports_in_file(file_path):
            print(f"Updated imports in {file_path}")
            updated_files += 1
    
    print(f"\nSummary: Updated imports in {updated_files} out of {len(python_files)} Python files.")


if __name__ == '__main__':
    main()
