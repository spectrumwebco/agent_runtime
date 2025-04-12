"""
Script to fix common linting issues in Python code.
"""

import os
import re
import sys
from pathlib import Path


def fix_unused_imports(file_path):
    """Remove unused imports marked by flake8 F401."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    import_pattern = r'^import\s+([^\n]+)$|^from\s+([^\s]+)\s+import\s+([^\n]+)$'
    imports = re.findall(import_pattern, content, re.MULTILINE)
    
    unused_imports = []
    for imp in imports:
        if imp[0]:  # Simple import
            module = imp[0].strip()
            if f"'{module}' imported but unused" in content:
                unused_imports.append(module)
        elif imp[1] and imp[2]:  # From import
            module = imp[1].strip()
            items = [item.strip() for item in imp[2].split(',')]
            for item in items:
                if f"'{item}' imported but unused" in content:
                    unused_imports.append(f"{module}.{item}")
    
    for unused in unused_imports:
        if '.' in unused:
            module, item = unused.rsplit('.', 1)
            pattern = fr'from\s+{module}\s+import\s+([^,\n]*{item}[^,\n]*)(,\s*|\n)'
            content = re.sub(pattern, r'\2', content)
            content = re.sub(fr'from\s+{module}\s+import\s*\n', '', content)
        else:
            content = re.sub(fr'import\s+{unused}\n', '', content)
    
    with open(file_path, 'w') as f:
        f.write(content)


def fix_long_lines(file_path):
    """Fix lines that are too long (E501)."""
    with open(file_path, 'r') as f:
        lines = f.readlines()
    
    modified_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]
        
        if len(line.rstrip('\n')) > 79:
            if ('(' in line and ')' not in line) or ('(' in line and line.count('(') > line.count(')')):
                new_line = line.rstrip('\n')
                j = i + 1
                
                while j < len(lines) and (')' not in lines[j] or lines[j].count('(') > lines[j].count(')')):
                    new_line += lines[j].rstrip('\n')
                    j += 1
                
                if j < len(lines):
                    new_line += lines[j].rstrip('\n')
                    j += 1
                
                if '(' in new_line and ')' in new_line:
                    open_idx = new_line.find('(')
                    close_idx = new_line.rfind(')')
                    
                    prefix = new_line[:open_idx+1]
                    content = new_line[open_idx+1:close_idx]
                    suffix = new_line[close_idx:]
                    
                    parts = content.split(',')
                    
                    formatted = prefix + '\n'
                    indent = ' ' * (len(prefix) - 1)
                    
                    for k, part in enumerate(parts):
                        if k < len(parts) - 1:
                            formatted += indent + part.strip() + ',\n'
                        else:
                            formatted += indent + part.strip() + '\n'
                    
                    formatted += indent[:-4] + suffix + '\n'
                    
                    modified_lines.append(formatted)
                    i = j
                    continue
            
            if len(line) > 79 and '=' in line:
                split_point = line.find('=') + 1
                modified_lines.append(line[:split_point] + '\n')
                modified_lines.append(' ' * 4 + line[split_point:].lstrip())
            else:
                modified_lines.append(line)
        else:
            modified_lines.append(line)
        
        i += 1
    
    with open(file_path, 'w') as f:
        f.writelines(modified_lines)


def main():
    """Main function to fix linting issues."""
    src_dir = Path('src')
    
    if not src_dir.exists():
        print(f"Error: {src_dir} directory not found")
        sys.exit(1)
    
    python_files = list(src_dir.glob('**/*.py'))
    
    for file_path in python_files:
        print(f"Processing {file_path}...")
        fix_unused_imports(file_path)
        fix_long_lines(file_path)
    
    print(f"Processed {len(python_files)} Python files")


if __name__ == '__main__':
    main()
