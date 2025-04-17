"""
Deployment Module for Rumble

This module specializes in setting up deployment infrastructure and configurations
with a focus on deployment excellence. It handles build processes, CI/CD pipelines,
environment configurations, and monitoring setup.
"""

import os
import json
from typing import Dict, List, Any, Optional

class BaseModule:
    """Base class for all modules in the agent system."""
    
    @property
    def name(self) -> str:
        """Returns the name of the module."""
        raise NotImplementedError("Subclasses must implement this method")
    
    @property
    def description(self) -> str:
        """Returns the description of the module."""
        raise NotImplementedError("Subclasses must implement this method")
    
    @property
    def tools(self) -> List[Dict[str, Any]]:
        """Returns the list of tools provided by this module."""
        raise NotImplementedError("Subclasses must implement this method")
    
    def initialize(self, context: Dict[str, Any]) -> None:
        """Initializes the module with the provided context."""
        raise NotImplementedError("Subclasses must implement this method")
    
    def cleanup(self) -> None:
        """Cleans up any resources used by the module."""
        raise NotImplementedError("Subclasses must implement this method")


class DeploymentModule(BaseModule):
    """
    Deployment Module for setting up deployment infrastructure and configurations.
    
    This module extends the Datasource Module with specialized capabilities for
    configuring production builds, creating CI/CD pipeline configurations, setting
    up environment variables, and implementing logging and monitoring.
    """
    
    def __init__(self):
        """Initialize the Deployment Module."""
        self._context = {}
        self._templates = {}
        self._deployment_registry = {}
    
    @property
    def name(self) -> str:
        """Returns the name of the module."""
        return "deployment"
    
    @property
    def description(self) -> str:
        """Returns the description of the module."""
        return "Handles setting up deployment infrastructure and configurations, including build process, CI/CD setup, environment configuration, and monitoring setup."
    
    @property
    def tools(self) -> List[Dict[str, Any]]:
        """Returns the list of tools provided by this module."""
        return [
            {
                "name": "configure_build_process",
                "description": "Configure production build process",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "project_type": {
                            "type": "string",
                            "description": "Type of project (e.g., frontend, backend, fullstack, etc.)"
                        },
                        "build_tool": {
                            "type": "string",
                            "description": "Build tool to use (e.g., webpack, vite, etc.)"
                        }
                    },
                    "required": ["project_dir", "project_type"]
                }
            },
            {
                "name": "setup_ci_cd",
                "description": "Create CI/CD pipeline configurations",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "repository_type": {
                            "type": "string",
                            "description": "Type of repository (e.g., github, gitlab, etc.)"
                        },
                        "ci_cd_platform": {
                            "type": "string",
                            "description": "CI/CD platform to use (e.g., github-actions, gitlab-ci, etc.)"
                        }
                    },
                    "required": ["project_dir", "repository_type", "ci_cd_platform"]
                }
            },
            {
                "name": "configure_environments",
                "description": "Set up environment-specific configurations",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "environments": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            },
                            "description": "Environments to configure (e.g., development, staging, production)"
                        },
                        "project_type": {
                            "type": "string",
                            "description": "Type of project (e.g., frontend, backend, fullstack, etc.)"
                        }
                    },
                    "required": ["project_dir", "environments", "project_type"]
                }
            },
            {
                "name": "setup_monitoring",
                "description": "Implement logging and monitoring",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "monitoring_type": {
                            "type": "string",
                            "description": "Type of monitoring to set up (e.g., logging, error-tracking, performance, etc.)"
                        },
                        "project_type": {
                            "type": "string",
                            "description": "Type of project (e.g., frontend, backend, fullstack, etc.)"
                        }
                    },
                    "required": ["project_dir", "monitoring_type", "project_type"]
                }
            }
        ]
    
    def initialize(self, context: Dict[str, Any]) -> None:
        """Initializes the module with the provided context."""
        self._context = context
        self._load_templates()
    
    def cleanup(self) -> None:
        """Cleans up any resources used by the module."""
        self._context = {}
        self._templates = {}
        self._deployment_registry = {}
    
    def _load_templates(self) -> None:
        """Load deployment templates from the templates directory."""
        templates_dir = os.path.join(os.path.dirname(__file__), "templates", "deployment")
        if not os.path.exists(templates_dir):
            os.makedirs(templates_dir)
        
        for root, dirs, files in os.walk(templates_dir):
            for file in files:
                if file.endswith(".json"):
                    template_path = os.path.join(root, file)
                    try:
                        with open(template_path, "r") as f:
                            template = json.load(f)
                            template_id = template.get("id")
                            if template_id:
                                self._templates[template_id] = template
                    except Exception as e:
                        print(f"Error loading template {template_path}: {e}")
    
    def configure_build_process(self, build_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Configure production build process.
        
        Args:
            build_spec: Dictionary containing build specifications
            
        Returns:
            Dictionary containing the configuration results
        """
        project_dir = build_spec.get("project_dir", "")
        project_type = build_spec.get("project_type", "")
        build_tool = build_spec.get("build_tool", "")
        
        if not all([project_dir, project_type]):
            return {"error": "Missing required parameters"}
        
        if project_type.lower() == "frontend":
            if build_tool.lower() == "webpack":
                return self._configure_webpack(project_dir)
            elif build_tool.lower() == "vite":
                return self._configure_vite(project_dir)
            elif not build_tool:
                if os.path.exists(os.path.join(project_dir, "vite.config.js")) or os.path.exists(os.path.join(project_dir, "vite.config.ts")):
                    return self._configure_vite(project_dir)
                elif os.path.exists(os.path.join(project_dir, "webpack.config.js")):
                    return self._configure_webpack(project_dir)
                else:
                    return {"error": "Could not auto-detect build tool, please specify build_tool parameter"}
            else:
                return {"error": f"Unsupported build tool: {build_tool}"}
        elif project_type.lower() == "backend":
            if project_dir.endswith("django"):
                return self._configure_django_build(project_dir)
            elif project_dir.endswith("flask"):
                return self._configure_flask_build(project_dir)
            elif project_dir.endswith("express"):
                return self._configure_express_build(project_dir)
            else:
                return {"error": f"Could not determine backend framework for {project_dir}"}
        else:
            return {"error": f"Unsupported project type: {project_type}"}
    
    def _configure_vite(self, project_dir: str) -> Dict[str, Any]:
        """
        Configure Vite build process for a frontend project.
        
        Args:
            project_dir: Directory of the project
            
        Returns:
            Dictionary containing the configuration results
        """
        files = {}
        
        vite_config_file = None
        if os.path.exists(os.path.join(project_dir, "vite.config.js")):
            vite_config_file = os.path.join(project_dir, "vite.config.js")
        elif os.path.exists(os.path.join(project_dir, "vite.config.ts")):
            vite_config_file = os.path.join(project_dir, "vite.config.ts")
        else:
            vite_config_file = os.path.join(project_dir, "vite.config.js")
            vite_config_content = """import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    minify: true,
    target: 'es2015',
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
        },
      },
    },
  },
});
"""
            files[vite_config_file] = vite_config_content
        
        env_production_file = os.path.join(project_dir, ".env.production")
        env_production_content = """# Production environment variables
VITE_API_URL=https://api.example.com
VITE_APP_ENV=production
"""
        files[env_production_file] = env_production_content
        
        env_staging_file = os.path.join(project_dir, ".env.staging")
        env_staging_content = """# Staging environment variables
VITE_API_URL=https://staging-api.example.com
VITE_APP_ENV=staging
"""
        files[env_staging_file] = env_staging_content
        
        package_json_file = os.path.join(project_dir, "package.json")
        if os.path.exists(package_json_file):
            try:
                with open(package_json_file, "r") as f:
                    package_json = json.load(f)
                
                if "scripts" not in package_json:
                    package_json["scripts"] = {}
                
                package_json["scripts"].update({
                    "build": "vite build",
                    "build:staging": "vite build --mode staging",
                    "build:production": "vite build --mode production",
                    "preview": "vite preview"
                })
                
                files[package_json_file] = json.dumps(package_json, indent=2)
            except Exception as e:
                return {"error": f"Error updating package.json: {e}"}
        
        for file_path, content in files.items():
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w") as f:
                f.write(content)
        
        self._deployment_registry[project_dir] = {
            "build_tool": "vite",
            "project_type": "frontend",
            "path": project_dir
        }
        
        return {
            "build_tool": "vite",
            "project_type": "frontend",
            "path": project_dir,
            "files": list(files.keys()),
            "message": "Vite build process has been configured successfully for frontend project."
        }
    
    def setup_ci_cd(self, ci_cd_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Create CI/CD pipeline configurations.
        
        Args:
            ci_cd_spec: Dictionary containing CI/CD specifications
            
        Returns:
            Dictionary containing the setup results
        """
        project_dir = ci_cd_spec.get("project_dir", "")
        repository_type = ci_cd_spec.get("repository_type", "")
        ci_cd_platform = ci_cd_spec.get("ci_cd_platform", "")
        
        if not all([project_dir, repository_type, ci_cd_platform]):
            return {"error": "Missing required parameters"}
        
        if repository_type.lower() == "github" and ci_cd_platform.lower() == "github-actions":
            return self._setup_github_actions(project_dir)
        elif repository_type.lower() == "gitlab" and ci_cd_platform.lower() == "gitlab-ci":
            return self._setup_gitlab_ci(project_dir)
        else:
            return {"error": f"Unsupported repository type or CI/CD platform: {repository_type}, {ci_cd_platform}"}
    
    def _setup_github_actions(self, project_dir: str) -> Dict[str, Any]:
        """
        Set up GitHub Actions CI/CD for a project.
        
        Args:
            project_dir: Directory of the project
            
        Returns:
            Dictionary containing the setup results
        """
        files = {}
        
        is_frontend = os.path.exists(os.path.join(project_dir, "package.json"))
        is_python = os.path.exists(os.path.join(project_dir, "requirements.txt")) or os.path.exists(os.path.join(project_dir, "pyproject.toml"))
        
        github_actions_dir = os.path.join(project_dir, ".github", "workflows")
        os.makedirs(github_actions_dir, exist_ok=True)
        
        if is_frontend:
            ci_workflow_file = os.path.join(github_actions_dir, "ci.yml")
            ci_workflow_content = """name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        node-version: [16.x, 18.x]

    steps:
    - uses: actions/checkout@v3
    - name: Use Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v3
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
    - name: Install dependencies
      run: npm ci
    - name: Lint
      run: npm run lint
    - name: Test
      run: npm test
    - name: Build
      run: npm run build
"""
            files[ci_workflow_file] = ci_workflow_content
            
            cd_workflow_file = os.path.join(github_actions_dir, "cd.yml")
            cd_workflow_content = """name: CD

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3
    - name: Use Node.js 18.x
      uses: actions/setup-node@v3
      with:
        node-version: 18.x
        cache: 'npm'
    - name: Install dependencies
      run: npm ci
    - name: Build
      run: npm run build
    - name: Deploy to production
      uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./dist
"""
            files[cd_workflow_file] = cd_workflow_content
        
        elif is_python:
            ci_workflow_file = os.path.join(github_actions_dir, "ci.yml")
            ci_workflow_content = """name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        python-version: [3.9, 3.10, 3.11]

    steps:
    - uses: actions/checkout@v3
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
        if [ -f requirements-dev.txt ]; then pip install -r requirements-dev.txt; fi
        if [ -f pyproject.toml ]; then pip install -e .; fi
    - name: Lint with flake8
      run: |
        pip install flake8
        flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics
    - name: Test with pytest
      run: |
        pip install pytest pytest-cov
        pytest --cov=./ --cov-report=xml
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.xml
        fail_ci_if_error: true
"""
            files[ci_workflow_file] = ci_workflow_content
            
            cd_workflow_file = os.path.join(github_actions_dir, "cd.yml")
            cd_workflow_content = """name: CD

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3
    - name: Set up Python 3.11
      uses: actions/setup-python@v4
      with:
        python-version: 3.11
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install build twine
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
    - name: Build and publish
      env:
        TWINE_USERNAME: ${{ secrets.PYPI_USERNAME }}
        TWINE_PASSWORD: ${{ secrets.PYPI_PASSWORD }}
      run: |
        python -m build
        twine upload dist/*
"""
            files[cd_workflow_file] = cd_workflow_content
        
        else:
            ci_workflow_file = os.path.join(github_actions_dir, "ci.yml")
            ci_workflow_content = """name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3
    - name: Build
      run: echo "Add your build steps here"
    - name: Test
      run: echo "Add your test steps here"
"""
            files[ci_workflow_file] = ci_workflow_content
        
        for file_path, content in files.items():
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w") as f:
                f.write(content)
        
        self._deployment_registry[project_dir] = {
            "ci_cd_platform": "github-actions",
            "repository_type": "github",
            "path": project_dir
        }
        
        return {
            "ci_cd_platform": "github-actions",
            "repository_type": "github",
            "path": project_dir,
            "files": list(files.keys()),
            "message": "GitHub Actions CI/CD has been set up successfully for the project."
        }
    
    def configure_environments(self, env_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Set up environment-specific configurations.
        
        Args:
            env_spec: Dictionary containing environment specifications
            
        Returns:
            Dictionary containing the configuration results
        """
        project_dir = env_spec.get("project_dir", "")
        environments = env_spec.get("environments", [])
        project_type = env_spec.get("project_type", "")
        
        if not all([project_dir, environments, project_type]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Environment configurations for {', '.join(environments)} would be set up here"}
    
    def setup_monitoring(self, monitoring_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Implement logging and monitoring.
        
        Args:
            monitoring_spec: Dictionary containing monitoring specifications
            
        Returns:
            Dictionary containing the setup results
        """
        project_dir = monitoring_spec.get("project_dir", "")
        monitoring_type = monitoring_spec.get("monitoring_type", "")
        project_type = monitoring_spec.get("project_type", "")
        
        if not all([project_dir, monitoring_type, project_type]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Monitoring of type {monitoring_type} would be set up here"}
