"""
Script to generate Python gRPC code from proto files.

This script uses the grpcio-tools package to generate Python code
from the proto files in the protos directory.
"""

import os
import sys
import subprocess
import logging
from pathlib import Path

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


def generate_grpc_code(proto_dir, output_dir):
    """
    Generate Python gRPC code from proto files.

    Args:
        proto_dir: Directory containing proto files
        output_dir: Directory to output generated code
    """
    os.makedirs(output_dir, exist_ok=True)

    proto_files = list(Path(proto_dir).glob('*.proto'))

    if not proto_files:
        logger.error(f"No proto files found in {proto_dir}")
        return False

    logger.info(
        f"Found {len(proto_files)} proto files: "
        f"{[f.name for f in proto_files]}")

    for proto_file in proto_files:
        logger.info(f"Generating code for {proto_file.name}")

        init_file = Path(output_dir) / '__init__.py'
        if not init_file.exists():
            init_file.touch()

        cmd = [
            'python', '-m', 'grpc_tools.protoc',
            f'--proto_path={proto_dir}',
            f'--python_out={output_dir}',
            f'--grpc_python_out={output_dir}',
            str(proto_file)
        ]

        try:
            subprocess.run(cmd, check=True)
            logger.info(f"Successfully generated code for {proto_file.name}")
        except subprocess.CalledProcessError as e:
            logger.error(f"Failed to generate code for {proto_file.name}: {e}")
            return False

    return True


def main():
    """Main function."""
    base_dir = Path(__file__).resolve().parent.parent

    proto_dir = base_dir / 'protos'
    output_dir = base_dir / 'api' / 'generated'

    try:
        import grpc_tools  # noqa: F401
    except ImportError:
        logger.info("Installing grpcio-tools...")
        subprocess.run([sys.executable, '-m', 'pip',
                       'install', 'grpcio-tools'], check=True)

    if generate_grpc_code(proto_dir, output_dir):
        logger.info(f"Successfully generated gRPC code in {output_dir}")
    else:
        logger.error("Failed to generate gRPC code")
        sys.exit(1)


if __name__ == "__main__":
    main()
