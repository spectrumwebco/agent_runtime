"""
Script to run tests after code reorganization.
"""

import os
import sys
import subprocess
from pathlib import Path


def run_tests(test_dir):
    """Run pytest on the specified test directory."""
    print(f"Running tests in {test_dir}...")
    
    try:
        result = subprocess.run(
            ["python", "-m", "pytest", test_dir, "-v"],
            capture_output=True,
            text=True,
            check=False
        )
        
        print("\nTest Results:")
        print(result.stdout)
        
        if result.returncode != 0:
            print("\nTest Failures:")
            print(result.stderr)
            return False
        
        return True
    except Exception as e:
        print(f"Error running tests: {str(e)}")
        return False


def main():
    """Main function."""
    repo_dir = sys.argv[1] if len(sys.argv) > 1 else '/home/ubuntu/repos/Fine-Tune'
    test_dir = os.path.join(repo_dir, 'tests')
    
    if not os.path.exists(test_dir):
        print(f"Error: Test directory {test_dir} does not exist.")
        sys.exit(1)
    
    success = run_tests(test_dir)
    
    if success:
        print("\nAll tests passed successfully!")
        sys.exit(0)
    else:
        print("\nSome tests failed. Please check the output above for details.")
        sys.exit(1)


if __name__ == '__main__':
    main()
