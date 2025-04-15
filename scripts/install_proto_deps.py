"""
Script to install protobuf dependencies for Python and Go.

This script installs the required dependencies for generating protobuf code
for Python and Go. It uses pip for Python dependencies and go get for Go dependencies.
"""

import os
import subprocess
import sys
from pathlib import Path


def main():
    """Install protobuf dependencies for Python and Go."""
    print("Installing Python protobuf dependencies...")
    python_deps = [
        "grpcio",
        "grpcio-tools",
        "protobuf",
    ]
    
    try:
        subprocess.run([sys.executable, "-m", "pip", "install", *python_deps], check=True)
        print("Python dependencies installed successfully")
    except subprocess.CalledProcessError as e:
        print(f"Error installing Python dependencies: {e}")
        return 1
    
    print("\nInstalling Go protobuf dependencies...")
    go_deps = [
        "google.golang.org/protobuf/cmd/protoc-gen-go",
        "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
    ]
    
    try:
        for dep in go_deps:
            subprocess.run(["go", "install", f"{dep}@latest"], check=True)
        print("Go dependencies installed successfully")
    except subprocess.CalledProcessError as e:
        print(f"Error installing Go dependencies: {e}")
        return 1
    
    print("\nChecking for protoc compiler...")
    try:
        result = subprocess.run(["protoc", "--version"], capture_output=True, text=True, check=True)
        print(f"Found protoc: {result.stdout.strip()}")
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("protoc not found. Please install it manually.")
        print("Visit https://github.com/protocolbuffers/protobuf/releases")
        return 1
    
    print("\nAll dependencies installed successfully")
    return 0


if __name__ == "__main__":
    sys.exit(main())
