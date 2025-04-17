"""
Repository Tracking Module for Corki Agent - Backup Agent for 100% Code Coverage

This module is responsible for discovering, tracking, and monitoring all company repositories.
It provides tools for repository discovery, structure analysis, language detection, and Git integration.
"""

import json
import logging
import os
import subprocess
from typing import Any, Dict, List, Optional, Set, Tuple, Union

import git
from git import Repo

from base_module import BaseModule

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class RepositoryTrackingModule(BaseModule):
    """
    Module for tracking and analyzing repositories.
    
    This module provides tools for discovering repositories, analyzing their structure,
    detecting languages, tracking dependencies, and integrating with Git.
    """
    
    def __init__(self, config: Optional[Dict[str, Any]] = None):
        """
        Initialize the repository tracking module.
        
        Args:
            config: Optional configuration dictionary for the module
        """
        super().__init__("repository_tracking", config)
        self.repositories = {}
        self.language_extensions = {
            "go": [".go"],
            "python": [".py"],
            "typescript": [".ts", ".tsx"],
            "javascript": [".js", ".jsx"],
            "java": [".java"],
            "rust": [".rs"],
            "cpp": [".cpp", ".cc", ".cxx", ".h", ".hpp"],
            "csharp": [".cs"],
            "php": [".php"]
        }
    
    def initialize(self) -> bool:
        """
        Initialize the repository tracking module and register its tools.
        
        Returns:
            bool: True if initialization was successful, False otherwise
        """
        self.register_tool(
            "discover_repositories",
            self.discover_repositories,
            "Discover repositories in a directory",
            [
                {
                    "name": "root_dir",
                    "type": "string",
                    "description": "Root directory to search for repositories"
                },
                {
                    "name": "recursive",
                    "type": "boolean",
                    "description": "Whether to search recursively",
                    "default": True
                }
            ],
            {
                "type": "array",
                "description": "List of discovered repository paths"
            }
        )
        
        self.register_tool(
            "analyze_repository",
            self.analyze_repository,
            "Analyze a repository's structure and languages",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                }
            ],
            {
                "type": "object",
                "description": "Repository analysis results"
            }
        )
        
        self.register_tool(
            "detect_languages",
            self.detect_languages,
            "Detect programming languages used in a repository",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                }
            ],
            {
                "type": "object",
                "description": "Language detection results"
            }
        )
        
        self.register_tool(
            "track_changes",
            self.track_changes,
            "Track changes in a repository since last check",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                },
                {
                    "name": "reference_commit",
                    "type": "string",
                    "description": "Reference commit to compare against",
                    "default": "HEAD~1"
                }
            ],
            {
                "type": "object",
                "description": "Changes since last check"
            }
        )
        
        self.register_tool(
            "get_repository_info",
            self.get_repository_info,
            "Get information about a repository",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                }
            ],
            {
                "type": "object",
                "description": "Repository information"
            }
        )
        
        self.initialized = True
        return True
    
    def cleanup(self) -> bool:
        """
        Clean up any resources used by the module.
        
        Returns:
            bool: True if cleanup was successful, False otherwise
        """
        self.repositories = {}
        self.initialized = False
        return True
    
    def discover_repositories(self, root_dir: str, recursive: bool = True) -> List[str]:
        """
        Discover Git repositories in a directory.
        
        Args:
            root_dir: Root directory to search for repositories
            recursive: Whether to search recursively
            
        Returns:
            List[str]: List of discovered repository paths
        """
        if not os.path.exists(root_dir):
            logger.error(f"Root directory {root_dir} does not exist")
            return []
        
        repo_paths = []
        
        if recursive:
            for root, dirs, _ in os.walk(root_dir):
                if '.git' in dirs:
                    repo_paths.append(root)
        else:
            dirs = [d for d in os.listdir(root_dir) if os.path.isdir(os.path.join(root_dir, d))]
            for d in dirs:
                if os.path.exists(os.path.join(root_dir, d, '.git')):
                    repo_paths.append(os.path.join(root_dir, d))
        
        for repo_path in repo_paths:
            if repo_path not in self.repositories:
                self.repositories[repo_path] = {
                    "path": repo_path,
                    "name": os.path.basename(repo_path),
                    "last_analyzed": None,
                    "languages": {},
                    "structure": {},
                    "last_commit": None
                }
        
        logger.info(f"Discovered {len(repo_paths)} repositories in {root_dir}")
        return repo_paths
    
    def analyze_repository(self, repo_path: str) -> Dict[str, Any]:
        """
        Analyze a repository's structure and languages.
        
        Args:
            repo_path: Path to the repository
            
        Returns:
            Dict[str, Any]: Repository analysis results
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {}
        
        if repo_path not in self.repositories:
            self.repositories[repo_path] = {
                "path": repo_path,
                "name": os.path.basename(repo_path),
                "last_analyzed": None,
                "languages": {},
                "structure": {},
                "last_commit": None
            }
        
        languages = self.detect_languages(repo_path)
        self.repositories[repo_path]["languages"] = languages
        
        structure = self._analyze_structure(repo_path)
        self.repositories[repo_path]["structure"] = structure
        
        repo_info = self.get_repository_info(repo_path)
        self.repositories[repo_path].update(repo_info)
        
        import time
        self.repositories[repo_path]["last_analyzed"] = time.time()
        
        logger.info(f"Analyzed repository {repo_path}")
        return self.repositories[repo_path]
    
    def detect_languages(self, repo_path: str) -> Dict[str, float]:
        """
        Detect programming languages used in a repository.
        
        Args:
            repo_path: Path to the repository
            
        Returns:
            Dict[str, float]: Language detection results with percentages
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {}
        
        language_counts = {}
        total_files = 0
        
        for root, _, files in os.walk(repo_path):
            if '.git' in root.split(os.path.sep):
                continue
            
            for file in files:
                file_path = os.path.join(root, file)
                _, ext = os.path.splitext(file)
                
                for language, extensions in self.language_extensions.items():
                    if ext in extensions:
                        language_counts[language] = language_counts.get(language, 0) + 1
                        total_files += 1
                        break
        
        language_percentages = {}
        if total_files > 0:
            for language, count in language_counts.items():
                language_percentages[language] = round((count / total_files) * 100, 2)
        
        logger.info(f"Detected languages in {repo_path}: {language_percentages}")
        return language_percentages
    
    def track_changes(self, repo_path: str, reference_commit: str = "HEAD~1") -> Dict[str, Any]:
        """
        Track changes in a repository since a reference commit.
        
        Args:
            repo_path: Path to the repository
            reference_commit: Reference commit to compare against
            
        Returns:
            Dict[str, Any]: Changes since the reference commit
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {}
        
        try:
            repo = Repo(repo_path)
            
            diff_index = repo.head.commit.diff(reference_commit)
            
            changes = {
                "added": [],
                "deleted": [],
                "modified": [],
                "renamed": []
            }
            
            for diff in diff_index:
                if diff.new_file:
                    changes["added"].append(diff.b_path)
                elif diff.deleted_file:
                    changes["deleted"].append(diff.a_path)
                elif diff.renamed:
                    changes["renamed"].append({
                        "old_path": diff.a_path,
                        "new_path": diff.b_path
                    })
                else:
                    changes["modified"].append(diff.a_path)
            
            if repo_path in self.repositories:
                self.repositories[repo_path]["last_commit"] = str(repo.head.commit)
            
            logger.info(f"Tracked changes in {repo_path} since {reference_commit}")
            return changes
        except Exception as e:
            logger.error(f"Error tracking changes in {repo_path}: {str(e)}")
            return {}
    
    def get_repository_info(self, repo_path: str) -> Dict[str, Any]:
        """
        Get information about a repository.
        
        Args:
            repo_path: Path to the repository
            
        Returns:
            Dict[str, Any]: Repository information
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {}
        
        try:
            repo = Repo(repo_path)
            
            remotes = {}
            for remote in repo.remotes:
                remotes[remote.name] = [url for url in remote.urls]
            
            branches = [str(branch) for branch in repo.branches]
            
            current_branch = str(repo.active_branch)
            
            last_commit = {
                "hash": str(repo.head.commit),
                "author": str(repo.head.commit.author),
                "message": repo.head.commit.message,
                "date": str(repo.head.commit.committed_datetime)
            }
            
            repo_info = {
                "remotes": remotes,
                "branches": branches,
                "current_branch": current_branch,
                "last_commit": last_commit
            }
            
            logger.info(f"Got repository info for {repo_path}")
            return repo_info
        except Exception as e:
            logger.error(f"Error getting repository info for {repo_path}: {str(e)}")
            return {}
    
    def _analyze_structure(self, repo_path: str) -> Dict[str, Any]:
        """
        Analyze the structure of a repository.
        
        Args:
            repo_path: Path to the repository
            
        Returns:
            Dict[str, Any]: Repository structure analysis
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {}
        
        structure = {
            "directories": {},
            "files_by_type": {},
            "total_files": 0,
            "total_directories": 0
        }
        
        for root, dirs, files in os.walk(repo_path):
            if '.git' in root.split(os.path.sep):
                continue
            
            rel_path = os.path.relpath(root, repo_path)
            if rel_path == '.':
                rel_path = ''
            
            if rel_path:
                parts = rel_path.split(os.path.sep)
                current = structure["directories"]
                for part in parts:
                    if part not in current:
                        current[part] = {"files": [], "directories": {}}
                    current = current[part]["directories"]
                structure["total_directories"] += 1
            
            for file in files:
                _, ext = os.path.splitext(file)
                if ext:
                    structure["files_by_type"][ext] = structure["files_by_type"].get(ext, 0) + 1
                
                if rel_path:
                    parts = rel_path.split(os.path.sep)
                    current = structure["directories"]
                    for part in parts:
                        current = current[part]["directories"]
                    current["files"].append(file)
                else:
                    if "files" not in structure["directories"]:
                        structure["directories"]["files"] = []
                    structure["directories"]["files"].append(file)
                
                structure["total_files"] += 1
        
        logger.info(f"Analyzed structure of {repo_path}: {structure['total_files']} files, {structure['total_directories']} directories")
        return structure
