"""
Testing Module for Rumble

This module specializes in setting up comprehensive testing infrastructure
with a focus on code quality and developer experience. It handles unit testing,
integration testing, and end-to-end testing setup.
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


class TestingModule(BaseModule):
    """
    Testing Module for setting up comprehensive testing infrastructure.
    
    This module extends the Datasource Module with specialized capabilities for
    setting up test frameworks, generating unit tests, configuring integration
    testing, and setting up end-to-end testing infrastructure.
    """
    
    def __init__(self):
        """Initialize the Testing Module."""
        self._context = {}
        self._templates = {}
        self._test_registry = {}
    
    @property
    def name(self) -> str:
        """Returns the name of the module."""
        return "testing"
    
    @property
    def description(self) -> str:
        """Returns the description of the module."""
        return "Handles setting up comprehensive testing infrastructure, including test framework setup, unit test generation, integration test setup, and E2E testing."
    
    @property
    def tools(self) -> List[Dict[str, Any]]:
        """Returns the list of tools provided by this module."""
        return [
            {
                "name": "setup_test_framework",
                "description": "Configure testing libraries and frameworks",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Framework to use (e.g., jest, pytest, etc.)"
                        },
                        "project_type": {
                            "type": "string",
                            "description": "Type of project (e.g., frontend, backend, fullstack, etc.)"
                        }
                    },
                    "required": ["project_dir", "framework", "project_type"]
                }
            },
            {
                "name": "generate_unit_tests",
                "description": "Create unit tests for components",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "component_path": {
                            "type": "string",
                            "description": "Path to the component to test"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Testing framework to use (e.g., jest, pytest, etc.)"
                        }
                    },
                    "required": ["project_dir", "component_path", "framework"]
                }
            },
            {
                "name": "setup_integration_tests",
                "description": "Configure integration testing",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "api_path": {
                            "type": "string",
                            "description": "Path to the API to test"
                        },
                        "framework": {
                            "type": "string",
                            "description": "Testing framework to use (e.g., jest, pytest, etc.)"
                        }
                    },
                    "required": ["project_dir", "api_path", "framework"]
                }
            },
            {
                "name": "configure_e2e_testing",
                "description": "Set up end-to-end testing infrastructure",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "project_dir": {
                            "type": "string",
                            "description": "Directory of the project"
                        },
                        "framework": {
                            "type": "string",
                            "description": "E2E testing framework to use (e.g., cypress, playwright, etc.)"
                        },
                        "project_type": {
                            "type": "string",
                            "description": "Type of project (e.g., frontend, backend, fullstack, etc.)"
                        }
                    },
                    "required": ["project_dir", "framework", "project_type"]
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
        self._test_registry = {}
    
    def _load_templates(self) -> None:
        """Load test templates from the templates directory."""
        templates_dir = os.path.join(os.path.dirname(__file__), "templates", "testing")
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
    
    def setup_test_framework(self, setup_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Configure testing libraries and frameworks.
        
        Args:
            setup_spec: Dictionary containing setup specifications
            
        Returns:
            Dictionary containing the setup results
        """
        project_dir = setup_spec.get("project_dir", "")
        framework = setup_spec.get("framework", "")
        project_type = setup_spec.get("project_type", "")
        
        if not all([project_dir, framework, project_type]):
            return {"error": "Missing required parameters"}
        
        if project_type.lower() == "frontend":
            if framework.lower() == "jest":
                return self._setup_jest(project_dir)
            elif framework.lower() == "vitest":
                return self._setup_vitest(project_dir)
            else:
                return {"error": f"Unsupported frontend testing framework: {framework}"}
        elif project_type.lower() == "backend":
            if framework.lower() == "pytest":
                return self._setup_pytest(project_dir)
            elif framework.lower() == "jest":
                return self._setup_jest_backend(project_dir)
            else:
                return {"error": f"Unsupported backend testing framework: {framework}"}
        else:
            return {"error": f"Unsupported project type: {project_type}"}
    
    def _setup_jest(self, project_dir: str) -> Dict[str, Any]:
        """
        Set up Jest testing framework for a frontend project.
        
        Args:
            project_dir: Directory of the project
            
        Returns:
            Dictionary containing the setup results
        """
        files = {}
        
        jest_config_file = os.path.join(project_dir, "jest.config.js")
        jest_config_content = """module.exports = {
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(js|jsx|ts|tsx)$': ['babel-jest', { presets: ['@babel/preset-env', '@babel/preset-react', '@babel/preset-typescript'] }],
  },
  moduleNameMapper: {
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.js'],
  testMatch: ['**/__tests__/**/*.[jt]s?(x)', '**/?(*.)+(spec|test).[jt]s?(x)'],
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.d.ts',
    '!src/index.{js,jsx,ts,tsx}',
  ],
  coverageThreshold: {
    global: {
      branches: 70,
      functions: 70,
      lines: 70,
      statements: 70,
    },
  },
};
"""
        files[jest_config_file] = jest_config_content
        
        setup_tests_file = os.path.join(project_dir, "src", "setupTests.js")
        setup_tests_content = """// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';
"""
        files[setup_tests_file] = setup_tests_content
        
        package_json_file = os.path.join(project_dir, "package.json")
        if os.path.exists(package_json_file):
            try:
                with open(package_json_file, "r") as f:
                    package_json = json.load(f)
                
                if "devDependencies" not in package_json:
                    package_json["devDependencies"] = {}
                
                package_json["devDependencies"].update({
                    "jest": "^29.5.0",
                    "@testing-library/jest-dom": "^5.16.5",
                    "@testing-library/react": "^14.0.0",
                    "@testing-library/user-event": "^14.4.3",
                    "babel-jest": "^29.5.0",
                    "@babel/preset-env": "^7.21.4",
                    "@babel/preset-react": "^7.18.6",
                    "@babel/preset-typescript": "^7.21.4",
                    "identity-obj-proxy": "^3.0.0",
                    "jest-environment-jsdom": "^29.5.0"
                })
                
                if "scripts" not in package_json:
                    package_json["scripts"] = {}
                
                package_json["scripts"].update({
                    "test": "jest",
                    "test:watch": "jest --watch",
                    "test:coverage": "jest --coverage"
                })
                
                files[package_json_file] = json.dumps(package_json, indent=2)
            except Exception as e:
                return {"error": f"Error updating package.json: {e}"}
        
        for file_path, content in files.items():
            os.makedirs(os.path.dirname(file_path), exist_ok=True)
            with open(file_path, "w") as f:
                f.write(content)
        
        self._test_registry[project_dir] = {
            "framework": "jest",
            "project_type": "frontend",
            "path": project_dir
        }
        
        return {
            "framework": "jest",
            "project_type": "frontend",
            "path": project_dir,
            "files": list(files.keys()),
            "message": "Jest testing framework has been set up successfully for frontend project."
        }
    
    def generate_unit_tests(self, test_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Create unit tests for components.
        
        Args:
            test_spec: Dictionary containing test specifications
            
        Returns:
            Dictionary containing the generated test files
        """
        project_dir = test_spec.get("project_dir", "")
        component_path = test_spec.get("component_path", "")
        framework = test_spec.get("framework", "")
        
        if not all([project_dir, component_path, framework]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Unit tests for component at {component_path} would be generated here"}
    
    def setup_integration_tests(self, integration_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Configure integration testing.
        
        Args:
            integration_spec: Dictionary containing integration test specifications
            
        Returns:
            Dictionary containing the setup results
        """
        project_dir = integration_spec.get("project_dir", "")
        api_path = integration_spec.get("api_path", "")
        framework = integration_spec.get("framework", "")
        
        if not all([project_dir, api_path, framework]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"Integration tests for API at {api_path} would be set up here"}
    
    def configure_e2e_testing(self, e2e_spec: Dict[str, Any]) -> Dict[str, Any]:
        """
        Set up end-to-end testing infrastructure.
        
        Args:
            e2e_spec: Dictionary containing E2E test specifications
            
        Returns:
            Dictionary containing the setup results
        """
        project_dir = e2e_spec.get("project_dir", "")
        framework = e2e_spec.get("framework", "")
        project_type = e2e_spec.get("project_type", "")
        
        if not all([project_dir, framework, project_type]):
            return {"error": "Missing required parameters"}
        
        return {"message": f"E2E testing with {framework} would be set up here"}
