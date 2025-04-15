"""
Script to update import paths in the backend/apps directory to reflect the new Django structure.
"""

import os
import re
from pathlib import Path

def update_imports(file_path):
    """Update import paths in a Python file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    updated_content = re.sub(
        r'from python_agent\.', 
        'from apps.python_agent.', 
        content
    )
    updated_content = re.sub(
        r'import python_agent\.', 
        'import apps.python_agent.', 
        updated_content
    )
    
    updated_content = re.sub(
        r'from src\.ml_infrastructure\.', 
        'from apps.python_ml.', 
        updated_content
    )
    updated_content = re.sub(
        r'import src\.ml_infrastructure\.', 
        'import apps.python_ml.', 
        updated_content
    )
    
    updated_content = re.sub(
        r'from tools\.', 
        'from apps.tools.', 
        updated_content
    )
    updated_content = re.sub(
        r'import tools\.', 
        'import apps.tools.', 
        updated_content
    )
    
    if content != updated_content:
        with open(file_path, 'w') as f:
            f.write(updated_content)
        return True
    return False

def main():
    """Main function to update imports in all Python files."""
    base_dir = Path(__file__).parent
    apps_dir = base_dir / 'backend' / 'apps'
    
    python_files = []
    for root, _, files in os.walk(apps_dir):
        for file in files:
            if file.endswith('.py'):
                python_files.append(os.path.join(root, file))
    
    updated_files = 0
    for file_path in python_files:
        if update_imports(file_path):
            updated_files += 1
            print(f"Updated imports in {file_path}")
    
    print(f"Updated imports in {updated_files} files")

if __name__ == "__main__":
    main()
