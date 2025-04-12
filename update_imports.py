import os
import re
import sys

def update_imports(directory):
    """Update import statements in Python files to reflect new directory structure."""
    print(f"Updating imports in {directory}...")
    
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.py'):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r') as f:
                        content = f.read()
                    
                    updated_content = re.sub(
                        r'from agent_framework', 
                        r'from python_agent.agent_framework', 
                        content
                    )
                    updated_content = re.sub(
                        r'import agent_framework', 
                        r'import python_agent.agent_framework', 
                        updated_content
                    )
                    updated_content = re.sub(
                        r'from rex', 
                        r'from python_agent.agent_framework.framework', 
                        updated_content
                    )
                    updated_content = re.sub(
                        r'import rex', 
                        r'import python_agent.agent_framework.framework', 
                        updated_content
                    )
                    
                    if content != updated_content:
                        with open(file_path, 'w') as f:
                            f.write(updated_content)
                        print(f"Updated imports in {file_path}")
                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
    
    print("Import updates completed!")

if __name__ == "__main__":
    update_imports("/home/ubuntu/repos/agent_runtime/python_agent")
