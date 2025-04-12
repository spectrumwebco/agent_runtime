
import sys
import yaml

def validate_yaml(file_path):
    try:
        with open(file_path, 'r') as file:
            list(yaml.safe_load_all(file))
        print(f"✅ YAML validation successful for {file_path}")
        return True
    except yaml.YAMLError as e:
        print(f"❌ YAML validation failed for {file_path}: {e}")
        return False
    except Exception as e:
        print(f"❌ Error reading file {file_path}: {e}")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python validate_yaml.py <file1.yaml> [file2.yaml ...]")
        sys.exit(1)
    
    success = True
    for file_path in sys.argv[1:]:
        if not validate_yaml(file_path):
            success = False
    
    sys.exit(0 if success else 1)
