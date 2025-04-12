import os
import shutil
import re
import sys

def main():
    """Reorganize agent_framework directory to python_agent/agent_framework."""
    print("Starting reorganization...")
    
    source_dir = "/home/ubuntu/repos/agent_runtime/agent_framework"
    target_dir = "/home/ubuntu/repos/agent_runtime/python_agent/agent_framework"
    
    os.makedirs(target_dir, exist_ok=True)
    print(f"Created directory: {target_dir}")
    
    for item in os.listdir(source_dir):
        source_item = os.path.join(source_dir, item)
        target_item = os.path.join(target_dir, item)
        
        if os.path.isdir(source_item):
            if os.path.exists(target_item):
                shutil.rmtree(target_item)
            shutil.copytree(source_item, target_item)
            print(f"Copied directory: {item}")
        else:
            if os.path.exists(target_item):
                os.remove(target_item)
            shutil.copy2(source_item, target_item)
            print(f"Copied file: {item}")
    
    init_path = os.path.join(target_dir, "__init__.py")
    with open(init_path, 'w') as f:
        f.write('__version__ = "1.2.1"\n\n')
        f.write('REMOTE_EXECUTABLE_NAME = "agent-framework-remote"\n')
        f.write('PACKAGE_NAME = "agent-framework"\n')
    print("Updated __init__.py")
    
    for root, _, files in os.walk(target_dir):
        for file in files:
            if file.endswith('.py'):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r') as f:
                        content = f.read()
                    
                    updated_content = content.replace('swerex', 'agent_framework')
                    updated_content = updated_content.replace('SWE-ReX', 'agent_framework')
                    
                    if content != updated_content:
                        with open(file_path, 'w') as f:
                            f.write(updated_content)
                        print(f"Updated references in: {file_path}")
                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
    
    rex_dir = os.path.join(target_dir, "rex")
    framework_dir = os.path.join(target_dir, "framework")
    if os.path.isdir(rex_dir):
        if os.path.exists(framework_dir):
            shutil.rmtree(framework_dir)
        shutil.move(rex_dir, framework_dir)
        print("Renamed 'rex' folder to 'framework'")
    
    swerex_tests_dir = "/home/ubuntu/repos/agent_runtime/swerex_tests"
    agent_framework_tests_dir = "/home/ubuntu/repos/agent_runtime/agent_framework_tests"
    if os.path.isdir(swerex_tests_dir):
        if os.path.exists(agent_framework_tests_dir):
            shutil.rmtree(agent_framework_tests_dir)
        shutil.move(swerex_tests_dir, agent_framework_tests_dir)
        print("Renamed 'swerex_tests' to 'agent_framework_tests'")
    
    print("Reorganization completed successfully!")

if __name__ == "__main__":
    main()
