"""
Utilities for workspace management.

This module provides utilities for creating and managing workspaces
for conversations, including file operations and workspace creation.
"""

import os
import uuid
import shutil
import logging
from typing import List, Dict, Any, Optional
from django.conf import settings

logger = logging.getLogger(__name__)


def create_workspace() -> str:
    """
    Create a new workspace directory.
    
    Returns:
        str: The relative path to the workspace directory.
    """
    workspace_id = str(uuid.uuid4())
    
    workspace_path = os.path.join(settings.WORKSPACES_DIR, workspace_id)
    os.makedirs(workspace_path, exist_ok=True)
    
    logger.info(f"Created workspace at {workspace_path}")
    
    return workspace_id


def get_workspace_path(workspace_id: str) -> str:
    """
    Get the absolute path to a workspace.
    
    Args:
        workspace_id: The ID of the workspace.
        
    Returns:
        str: The absolute path to the workspace directory.
    """
    return os.path.join(settings.WORKSPACES_DIR, workspace_id)


def list_workspace_files(workspace_path: str) -> List[Dict[str, Any]]:
    """
    List all files in a workspace.
    
    Args:
        workspace_path: The absolute path to the workspace directory.
        
    Returns:
        List[Dict[str, Any]]: A list of file information dictionaries.
    """
    files = []
    
    for root, dirs, filenames in os.walk(workspace_path):
        for filename in filenames:
            file_path = os.path.join(root, filename)
            rel_path = os.path.relpath(file_path, workspace_path)
            
            if filename.startswith('.'):
                continue
            
            file_info = {
                'path': rel_path,
                'name': filename,
                'size': os.path.getsize(file_path),
                'last_modified': os.path.getmtime(file_path),
                'is_directory': False,
            }
            
            files.append(file_info)
    
    for root, dirs, _ in os.walk(workspace_path):
        for dirname in dirs:
            dir_path = os.path.join(root, dirname)
            rel_path = os.path.relpath(dir_path, workspace_path)
            
            if dirname.startswith('.'):
                continue
            
            dir_info = {
                'path': rel_path,
                'name': dirname,
                'size': 0,
                'last_modified': os.path.getmtime(dir_path),
                'is_directory': True,
            }
            
            files.append(dir_info)
    
    return files


def read_file_content(file_path: str) -> str:
    """
    Read the content of a file.
    
    Args:
        file_path: The absolute path to the file.
        
    Returns:
        str: The content of the file.
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()
    except UnicodeDecodeError:
        try:
            with open(file_path, 'r', encoding='latin-1') as f:
                return f.read()
        except Exception as e:
            logger.error(f"Error reading file {file_path}: {e}")
            return f"Error reading file: {e}"
    except Exception as e:
        logger.error(f"Error reading file {file_path}: {e}")
        return f"Error reading file: {e}"


def write_file_content(file_path: str, content: str) -> bool:
    """
    Write content to a file.
    
    Args:
        file_path: The absolute path to the file.
        content: The content to write.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        
        return True
    except Exception as e:
        logger.error(f"Error writing to file {file_path}: {e}")
        return False


def delete_file(file_path: str) -> bool:
    """
    Delete a file.
    
    Args:
        file_path: The absolute path to the file.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        if os.path.isfile(file_path):
            os.remove(file_path)
        elif os.path.isdir(file_path):
            shutil.rmtree(file_path)
        else:
            return False
        
        return True
    except Exception as e:
        logger.error(f"Error deleting file {file_path}: {e}")
        return False


def create_directory(dir_path: str) -> bool:
    """
    Create a directory.
    
    Args:
        dir_path: The absolute path to the directory.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        os.makedirs(dir_path, exist_ok=True)
        return True
    except Exception as e:
        logger.error(f"Error creating directory {dir_path}: {e}")
        return False


def copy_file(src_path: str, dst_path: str) -> bool:
    """
    Copy a file.
    
    Args:
        src_path: The absolute path to the source file.
        dst_path: The absolute path to the destination file.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        os.makedirs(os.path.dirname(dst_path), exist_ok=True)
        
        if os.path.isfile(src_path):
            shutil.copy2(src_path, dst_path)
        elif os.path.isdir(src_path):
            shutil.copytree(src_path, dst_path)
        else:
            return False
        
        return True
    except Exception as e:
        logger.error(f"Error copying file {src_path} to {dst_path}: {e}")
        return False


def move_file(src_path: str, dst_path: str) -> bool:
    """
    Move a file.
    
    Args:
        src_path: The absolute path to the source file.
        dst_path: The absolute path to the destination file.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        os.makedirs(os.path.dirname(dst_path), exist_ok=True)
        
        shutil.move(src_path, dst_path)
        return True
    except Exception as e:
        logger.error(f"Error moving file {src_path} to {dst_path}: {e}")
        return False


def rename_file(file_path: str, new_name: str) -> bool:
    """
    Rename a file.
    
    Args:
        file_path: The absolute path to the file.
        new_name: The new name for the file.
        
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        dir_path = os.path.dirname(file_path)
        new_path = os.path.join(dir_path, new_name)
        
        os.rename(file_path, new_path)
        return True
    except Exception as e:
        logger.error(f"Error renaming file {file_path} to {new_name}: {e}")
        return False


def get_file_info(file_path: str) -> Optional[Dict[str, Any]]:
    """
    Get information about a file.
    
    Args:
        file_path: The absolute path to the file.
        
    Returns:
        Optional[Dict[str, Any]]: A dictionary with file information, or None if the file doesn't exist.
    """
    try:
        if not os.path.exists(file_path):
            return None
        
        filename = os.path.basename(file_path)
        
        file_info = {
            'path': file_path,
            'name': filename,
            'size': os.path.getsize(file_path) if os.path.isfile(file_path) else 0,
            'last_modified': os.path.getmtime(file_path),
            'is_directory': os.path.isdir(file_path),
        }
        
        return file_info
    except Exception as e:
        logger.error(f"Error getting file info for {file_path}: {e}")
        return None
