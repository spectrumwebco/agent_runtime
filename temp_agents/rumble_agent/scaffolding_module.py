"""
Scaffolding Module for Rumble

This module specializes in project initialization and structure creation
with a focus on excellence in scaffolding new projects and adding features
to existing ones.
"""

import os
import json
import shutil
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


class ScaffoldingModule(BaseModule):
    """
    Scaffolding Module for project initialization and structure creation.
    
    This module extends the Planner Module with specialized capabilities for
    project initialization, template management, code generation, and
    configuration setup.
    """
    
    def __init__(self):
        """Initialize the Scaffolding Module."""
        self._context = {}
        self._templates = {}
        self._project_registry = {}
    
    @property
    def name(self) -> str:
        """Returns the name of the module."""
        return "scaffolding"
    
    @property
    def description(self) -> str:
        """Returns the description of the module."""
        return "Handles project initialization and structure creation, including template management, code generation, configuration setup, and documentation generation."
    
    @property
    def tools(self) -> List[Dict[str, Any]]:
        """Returns the list of tools provided by this module."""
        return [
            {
                "name": "scaffold_init_project",
                "description": "Initialize a new project with the specified framework",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_name": {
                            "type": "string",
                            "description": "Name of the project"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., react, vue, django, flask, etc.)"
                        },
                        "project_dir": {
                            "type": "string",
                            "description": "Directory where the project should be created"
                        },
                        "options": {
                            "type": "object",
                            "description": "Additional options for project initialization"
                        }
                    },
                    "required": ["project_name", "framework", "project_dir"]
                }
            },
            {
                "name": "scaffold_add_feature",
                "description": "Add a specific feature to an existing project",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "feature_name": {
                            "type": "string",
                            "description": "Name of the feature to add"
                        },
                        "feature_type": {
                            "type": "string",
                            "description": "Type of feature (e.g., component, page, api, etc.)"
                        },
                        "options": {
                            "type": "object",
                            "description": "Additional options for feature creation"
                        }
                    },
                    "required": ["project_dir", "feature_name", "feature_type"]
                }
            },
            {
                "name": "scaffold_setup_environment",
                "description": "Set up the development environment for a project",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "environment_type": {
                            "type": "string",
                            "description": "Type of environment to set up (e.g., development, testing, production)"
                        },
                        "options": {
                            "type": "object",
                            "description": "Additional options for environment setup"
                        }
                    },
                    "required": ["project_dir", "environment_type"]
                }
            },
            {
                "name": "scaffold_generate_documentation",
                "description": "Generate documentation for a project",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "documentation_type": {
                            "type": "string",
                            "description": "Type of documentation to generate (e.g., readme, api, user, developer)"
                        },
                        "options": {
                            "type": "object",
                            "description": "Additional options for documentation generation"
                        }
                    },
                    "required": ["project_dir", "documentation_type"]
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
        self._project_registry = {}
    
    def _load_templates(self) -> None:
        """Load project templates from the templates directory."""
        templates_dir = os.path.join(os.path.dirname(__file__), "templates", "scaffolding")
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
    
    def scaffold_init_project(self, init_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Initialize a new project with the specified framework.
        
        Args:
            init_spec: Dictionary containing initialization specifications
            
        Returns:
            Dictionary containing the initialization results
        """
        project_name = init_spec.get("project_name", "")
        framework = init_spec.get("framework", "")
        project_dir = init_spec.get("project_dir", "")
        options = init_spec.get("options", {})
        
        if not all([project_name, framework, project_dir]):
            return {"error": "Missing required parameters"}
        
        project_path = os.path.join(project_dir, project_name)
        if os.path.exists(project_path):
            return {"error": f"Project directory {project_path} already exists"}
        
        os.makedirs(project_path, exist_ok=True)
        
        if framework.lower() in ["react", "react-ts", "react-typescript"]:
            return self._init_react_project(project_name, project_path, options)
        elif framework.lower() in ["vue", "vue3", "vuejs"]:
            return self._init_vue_project(project_name, project_path, options)
        elif framework.lower() in ["django", "django-rest-framework", "drf"]:
            return self._init_django_project(project_name, project_path, options)
        elif framework.lower() in ["flask", "flask-restful"]:
            return self._init_flask_project(project_name, project_path, options)
        elif framework.lower() in ["express", "node", "nodejs"]:
            return self._init_express_project(project_name, project_path, options)
        else:
            return {"error": f"Unsupported framework: {framework}"}
    
    def _init_react_project(self, project_name: str, project_path: str, options: Dict[str, Any]) -> Dict[str, Any]:
        """
        Initialize a React project.
        
        Args:
            project_name: Name of the project
            project_path: Path to the project directory
            options: Additional options for project initialization
            
        Returns:
            Dictionary containing the initialization results
        """
        files = {}
        
        use_typescript = options.get("typescript", True)
        
        package_json = {
            "name": project_name,
            "version": "0.1.0",
            "private": True,
            "dependencies": {
                "react": "^18.2.0",
                "react-dom": "^18.2.0",
                "react-router-dom": "^6.10.0"
            },
            "devDependencies": {
                "vite": "^4.2.1",
                "@vitejs/plugin-react": "^3.1.0"
            },
            "scripts": {
                "dev": "vite",
                "build": "vite build",
                "preview": "vite preview"
            }
        }
        
        if use_typescript:
            package_json["devDependencies"].update({
                "typescript": "^5.0.4",
                "@types/react": "^18.0.28",
                "@types/react-dom": "^18.0.11"
            })
        
        package_json_file = os.path.join(project_path, "package.json")
        files[package_json_file] = json.dumps(package_json, indent=2)
        
        for file_path, content in files.items():
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w") as f:
                f.write(content)
        
        self._project_registry[project_path] = {
            "name": project_name,
            "framework": "react",
            "typescript": use_typescript,
            "path": project_path
        }
        
        return {
            "project_name": project_name,
            "framework": "react",
            "typescript": use_typescript,
            "path": project_path,
            "files": list(files.keys()),
            "message": f"React{'TypeScript' if use_typescript else ''} project {project_name} has been initialized successfully."
        }
