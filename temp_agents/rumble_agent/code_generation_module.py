"""
Code Generation Module for Rumble

This module specializes in creating high-quality initial code implementations
with a focus on code quality and developer experience. It handles component generation,
API integration, database setup, and authentication flows.
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


class CodeGenerationModule(BaseModule):
    """
    Code Generation Module for creating high-quality initial code implementations.
    
    This module extends the Knowledge Module with specialized capabilities for
    creating UI components, implementing API integration, setting up database
    connections, and implementing authentication flows.
    """
    
    def __init__(self):
        """Initialize the Code Generation Module."""
        self._context = {}
        self._templates = {}
        self._component_registry = {}
    
    @property
    def name(self) -> str:
        """Returns the name of the module."""
        return "code_generation"
    
    @property
    def description(self) -> str:
        """Returns the description of the module."""
        return "Handles creating high-quality initial code implementations, including UI components, API integration, database layer, and authentication."
    
    @property
    def tools(self) -> List[Dict[str, Any]]:
        """Returns the list of tools provided by this module."""
        return [
            {
                "name": "generate_component",
                "description": "Create UI components based on specifications",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "component_name": {
                            "type": "string",
                            "description": "Name of the component to generate"
                        },
                        "component_type": {
                            "type": "string",
                            "description": "Type of component (e.g., page, form, card, etc.)"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., react, vue, angular, etc.)"
                        },
                        "props": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            },
                            "description": "Properties of the component"
                        },
                        "state": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            },
                            "description": "State variables of the component"
                        }
                    },
                    "required": ["project_dir", "component_name", "component_type", "framework"]
                }
            },
            {
                "name": "generate_api_layer",
                "description": "Implement API integration code",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "api_name": {
                            "type": "string",
                            "description": "Name of the API to generate"
                        },
                        "endpoints": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            },
                            "description": "Endpoints of the API"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., react, vue, angular, etc.)"
                        }
                    },
                    "required": ["project_dir", "api_name", "endpoints", "framework"]
                }
            },
            {
                "name": "generate_database_layer",
                "description": "Set up database connections and models",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "database_type": {
                            "type": "string",
                            "description": "Type of database (e.g., postgresql, mongodb, etc.)"
                        },
                        "models": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            },
                            "description": "Models to generate"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., django, flask, express, etc.)"
                        }
                    },
                    "required": ["project_dir", "database_type", "models", "framework"]
                }
            },
            {
                "name": "generate_authentication",
                "description": "Implement authentication flows",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "auth_type": {
                            "type": "string",
                            "description": "Type of authentication (e.g., jwt, oauth, etc.)"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., react, vue, django, etc.)"
                        }
                    },
                    "required": ["project_dir", "auth_type", "framework"]
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
        self._component_registry = {}
    
    def _load_templates(self) -> None:
        """Load component templates from the templates directory."""
        templates_dir = os.path.join(os.path.dirname(__file__), "templates", "code_generation")
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
    
    def generate_component(self, component_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Create UI components based on specifications.
        
        Args:
            component_spec: Dictionary containing component specifications
            
        Returns:
            Dictionary containing the generated component files
        """
        project_dir = component_spec.get("project_dir", "")
        component_name = component_spec.get("component_name", "")
        component_type = component_spec.get("component_type", "")
        framework = component_spec.get("framework", "")
        props = component_spec.get("props", [])
        state = component_spec.get("state", [])
        
        if not all([project_dir, component_name, component_type, framework]):
            return {"error": "Missing required parameters"}
        
        if framework.lower() == "react":
            return self._generate_react_component(project_dir, component_name, component_type, props, state)
        elif framework.lower() == "vue":
            return self._generate_vue_component(project_dir, component_name, component_type, props, state)
        else:
            return {"error": f"Unsupported framework: {framework}"}
    
    def _generate_react_component(self, project_dir: str, component_name: str, component_type: str, 
                                 props: List[Dict[str, Any]], state: List[Dict[str, Any]]) -> Dict[str, Any]:
        """
        Generate a React component.
        
        Args:
            project_dir: Directory of the project
            component_name: Name of the component to generate
            component_type: Type of component (e.g., page, form, card, etc.)
            props: Properties of the component
            state: State variables of the component
            
        Returns:
            Dictionary containing the generated component files
        """
        files = {}
        
        use_typescript = self._is_typescript_project(project_dir)
        
        components_dir = os.path.join(project_dir, "src", "components")
        if component_type.lower() == "page":
            components_dir = os.path.join(project_dir, "src", "pages")
        
        component_dir = os.path.join(components_dir, component_name)
        os.makedirs(component_dir, exist_ok=True)
        
        component_file = os.path.join(component_dir, f"index.{'tsx' if use_typescript else 'jsx'}")
        
        component_content = f"""import React from 'react';
import './{component_name.lower()}.module.css';

{f"interface {component_name}Props {{}}" if use_typescript else ""}

const {component_name}{': React.FC<' + component_name + 'Props>' if use_typescript else ''} = () => {{
  return (
    <div className="{component_name.lower()}">
      <h2>{component_name}</h2>
      <p>This is a {component_name} component.</p>
    </div>
  );
}};

export default {component_name};
"""
        
        files[component_file] = component_content
        
        styles_file = os.path.join(component_dir, f"{component_name.lower()}.module.css")
        styles_content = f".{component_name.lower()} {{\n  /* Add your styles here */\n}}\n"
        files[styles_file] = styles_content
        
        for file_path, content in files.items():
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w") as f:
                f.write(content)
        
        self._component_registry[component_name] = {
            "name": component_name,
            "type": component_type,
            "framework": "react",
            "typescript": use_typescript,
            "path": component_dir
        }
        
        return {
            "component_name": component_name,
            "component_type": component_type,
            "framework": "react",
            "typescript": use_typescript,
            "files": list(files.keys()),
            "message": f"React{'TypeScript' if use_typescript else ''} {component_type} component {component_name} has been generated successfully."
        }
    
    def _is_typescript_project(self, project_dir: str) -> bool:
        """
        Determine if a project uses TypeScript.
        
        Args:
            project_dir: Directory of the project
            
        Returns:
            Boolean indicating if the project uses TypeScript
        """
        if os.path.exists(os.path.join(project_dir, "tsconfig.json")):
            return True
        
        for root, dirs, files in os.walk(project_dir):
            for file in files:
                if file.endswith(".ts") or file.endswith(".tsx"):
                    return True
        
        return False
    
    def generate_api_layer(self, api_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Implement API integration code.
        
        Args:
            api_spec: Dictionary containing API specifications
            
        Returns:
            Dictionary containing the generated API files
        """
        project_dir = api_spec.get("project_dir", "")
        api_name = api_spec.get("api_name", "")
        endpoints = api_spec.get("endpoints", [])
        framework = api_spec.get("framework", "")
        
        if not all([project_dir, api_name, endpoints, framework]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"API layer for {api_name} would be generated here"}
    
    def generate_database_layer(self, db_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Set up database connections and models.
        
        Args:
            db_spec: Dictionary containing database specifications
            
        Returns:
            Dictionary containing the generated database files
        """
        project_dir = db_spec.get("project_dir", "")
        database_type = db_spec.get("database_type", "")
        models = db_spec.get("models", [])
        framework = db_spec.get("framework", "")
        
        if not all([project_dir, database_type, models, framework]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Database layer for {database_type} would be generated here"}
    
    def generate_authentication(self, auth_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Implement authentication flows.
        
        Args:
            auth_spec: Dictionary containing authentication specifications
            
        Returns:
            Dictionary containing the generated authentication files
        """
        project_dir = auth_spec.get("project_dir", "")
        auth_type = auth_spec.get("auth_type", "")
        framework = auth_spec.get("framework", "")
        
        if not all([project_dir, auth_type, framework]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Authentication flow for {auth_type} would be generated here"}
