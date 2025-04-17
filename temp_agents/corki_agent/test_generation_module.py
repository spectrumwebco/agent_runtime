"""
Test Generation Module for Corki Agent - Backup Agent for 100% Code Coverage

This module is responsible for automatically generating tests for uncovered code.
It provides tools for test case generation, test template management, and test prioritization.
"""

import json
import logging
import os
from typing import Any, Dict, List, Optional, Set, Tuple, Union

from base_module import BaseModule

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class TestGenerationModule(BaseModule):
    """
    Module for generating tests for uncovered code.
    
    This module provides tools for generating test cases based on code analysis,
    managing test templates, and prioritizing tests based on coverage impact.
    """
    
    def __init__(self, config: Optional[Dict[str, Any]] = None):
        """
        Initialize the test generation module.
        
        Args:
            config: Optional configuration dictionary for the module
        """
        super().__init__("test_generation", config)
        self.template_dir = os.path.join(os.path.dirname(__file__), "templates")
        self.language_test_frameworks = {
            "go": "testing",
            "python": "pytest",
            "typescript": "jest",
            "java": "junit",
            "rust": "cargo-test",
            "cpp": "gtest",
            "csharp": "xunit",
            "php": "phpunit"
        }
    
    def initialize(self) -> bool:
        """
        Initialize the test generation module and register its tools.
        
        Returns:
            bool: True if initialization was successful, False otherwise
        """
        # Register tools
        self.register_tool(
            "generate_test",
            self.generate_test,
            "Generate a test for a specific file or function",
            [
                {
                    "name": "file_path",
                    "type": "string",
                    "description": "Path to the file to generate a test for"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language of the file"
                },
                {
                    "name": "function_name",
                    "type": "string",
                    "description": "Name of the function to generate a test for",
                    "required": False
                },
                {
                    "name": "test_framework",
                    "type": "string",
                    "description": "Test framework to use",
                    "required": False
                }
            ],
            {
                "type": "object",
                "description": "Generated test"
            }
        )
        
        self.register_tool(
            "generate_tests_for_uncovered_code",
            self.generate_tests_for_uncovered_code,
            "Generate tests for uncovered code in a repository",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language to generate tests for"
                },
                {
                    "name": "coverage_data",
                    "type": "object",
                    "description": "Coverage data from the code coverage module"
                }
            ],
            {
                "type": "object",
                "description": "Generated tests"
            }
        )
        
        self.register_tool(
            "get_test_template",
            self.get_test_template,
            "Get a test template for a specific language and framework",
            [
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
                },
                {
                    "name": "test_framework",
                    "type": "string",
                    "description": "Test framework",
                    "required": False
                }
            ],
            {
                "type": "string",
                "description": "Test template"
            }
        )
        
        self.register_tool(
            "prioritize_tests",
            self.prioritize_tests,
            "Prioritize tests based on coverage impact",
            [
                {
                    "name": "repo_path",
                    "type": "string",
                    "description": "Path to the repository"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
                },
                {
                    "name": "coverage_data",
                    "type": "object",
                    "description": "Coverage data from the code coverage module"
                }
            ],
            {
                "type": "array",
                "description": "Prioritized list of files to test"
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
        self.initialized = False
        return True
    
    def generate_test(self, file_path: str, language: str, function_name: Optional[str] = None, 
                     test_framework: Optional[str] = None) -> Dict[str, Any]:
        """
        Generate a test for a specific file or function.
        
        Args:
            file_path: Path to the file to generate a test for
            language: Programming language of the file
            function_name: Name of the function to generate a test for
            test_framework: Test framework to use
            
        Returns:
            Dict[str, Any]: Generated test
        """
        if not os.path.exists(file_path):
            logger.error(f"File {file_path} does not exist")
            return {"error": f"File {file_path} does not exist"}
        
        if language not in self.language_test_frameworks:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        # Use default test framework if not specified
        if not test_framework:
            test_framework = self.language_test_frameworks[language]
        
        # Read the file content
        try:
            with open(file_path, 'r') as f:
                file_content = f.read()
        except Exception as e:
            logger.error(f"Error reading file {file_path}: {str(e)}")
            return {"error": f"Error reading file {file_path}: {str(e)}"}
        
        # Get the test template
        template = self.get_test_template(language, test_framework)
        
        # Generate the test file path
        test_file_path = self._generate_test_file_path(file_path, language)
        
        # Parse the file to extract functions/methods
        functions = self._extract_functions(file_content, language, function_name)
        
        if not functions:
            logger.warning(f"No functions found in {file_path}")
            return {"error": f"No functions found in {file_path}"}
        
        # Generate test cases for each function
        test_cases = []
        for func in functions:
            test_case = self._generate_test_case(func, language, test_framework)
            test_cases.append(test_case)
        
        # Fill in the template
        import_path = os.path.splitext(os.path.basename(file_path))[0]
        test_content = self._fill_template(template, {
            "import_path": import_path,
            "test_cases": "\n\n".join(test_cases),
            "file_name": os.path.basename(file_path),
            "language": language,
            "test_framework": test_framework
        })
        
        result = {
            "test_file_path": test_file_path,
            "test_content": test_content,
            "functions_tested": [func["name"] for func in functions],
            "language": language,
            "test_framework": test_framework
        }
        
        logger.info(f"Generated test for {file_path} with {len(functions)} functions")
        return result
    
    def generate_tests_for_uncovered_code(self, repo_path: str, language: str, 
                                         coverage_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        Generate tests for uncovered code in a repository.
        
        Args:
            repo_path: Path to the repository
            language: Programming language to generate tests for
            coverage_data: Coverage data from the code coverage module
            
        Returns:
            Dict[str, Any]: Generated tests
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return {"error": f"Repository path {repo_path} does not exist"}
        
        if language not in self.language_test_frameworks:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        if "error" in coverage_data:
            logger.error(f"Invalid coverage data: {coverage_data['error']}")
            return {"error": f"Invalid coverage data: {coverage_data['error']}"}
        
        # Get files with low coverage
        files = coverage_data.get("files", {})
        uncovered_files = []
        
        for file_path, file_data in files.items():
            coverage = file_data.get("coverage", 0)
            
            # Consider files with less than 80% coverage as uncovered
            if coverage < 80:
                uncovered_files.append({
                    "file_path": os.path.join(repo_path, file_path),
                    "coverage": coverage,
                    "uncovered_lines": file_data.get("uncovered_lines", [])
                })
        
        # Sort files by coverage (lowest first)
        uncovered_files.sort(key=lambda x: x["coverage"])
        
        # Generate tests for each file
        generated_tests = []
        
        for file_data in uncovered_files:
            file_path = file_data["file_path"]
            
            # Skip test files
            if self._is_test_file(file_path, language):
                continue
            
            # Generate test
            test_result = self.generate_test(file_path, language)
            
            if "error" not in test_result:
                generated_tests.append({
                    "file_path": file_path,
                    "test_file_path": test_result["test_file_path"],
                    "test_content": test_result["test_content"],
                    "coverage_before": file_data["coverage"]
                })
        
        result = {
            "total_files": len(uncovered_files),
            "generated_tests": len(generated_tests),
            "tests": generated_tests
        }
        
        logger.info(f"Generated {len(generated_tests)} tests for uncovered code in {repo_path}")
        return result
    
    def get_test_template(self, language: str, test_framework: Optional[str] = None) -> str:
        """
        Get a test template for a specific language and framework.
        
        Args:
            language: Programming language
            test_framework: Test framework
            
        Returns:
            str: Test template
        """
        if language not in self.language_test_frameworks:
            logger.error(f"Unsupported language: {language}")
            return f"Error: Unsupported language: {language}"
        
        # Use default test framework if not specified
        if not test_framework:
            test_framework = self.language_test_frameworks[language]
        
        # Get template file path
        template_file = os.path.join(self.template_dir, language, f"{test_framework}_template.txt")
        
        # Check if template file exists
        if not os.path.exists(template_file):
            # Use default template
            return self._get_default_template(language, test_framework)
        
        # Read template file
        try:
            with open(template_file, 'r') as f:
                template = f.read()
            
            logger.info(f"Loaded test template for {language}/{test_framework}")
            return template
        except Exception as e:
            logger.error(f"Error reading template file {template_file}: {str(e)}")
            return self._get_default_template(language, test_framework)
    
    def prioritize_tests(self, repo_path: str, language: str, 
                        coverage_data: Dict[str, Any]) -> List[Dict[str, Any]]:
        """
        Prioritize tests based on coverage impact.
        
        Args:
            repo_path: Path to the repository
            language: Programming language
            coverage_data: Coverage data from the code coverage module
            
        Returns:
            List[Dict[str, Any]]: Prioritized list of files to test
        """
        if not os.path.exists(repo_path):
            logger.error(f"Repository path {repo_path} does not exist")
            return []
        
        if "error" in coverage_data:
            logger.error(f"Invalid coverage data: {coverage_data['error']}")
            return []
        
        # Get files with low coverage
        files = coverage_data.get("files", {})
        uncovered_files = []
        
        for file_path, file_data in files.items():
            coverage = file_data.get("coverage", 0)
            statements = file_data.get("statements", 0)
            covered = file_data.get("covered", 0)
            
            # Calculate potential impact
            uncovered_statements = statements - covered
            impact = uncovered_statements * (1 - (coverage / 100))
            
            # Consider files with less than 80% coverage as uncovered
            if coverage < 80:
                uncovered_files.append({
                    "file_path": os.path.join(repo_path, file_path),
                    "coverage": coverage,
                    "uncovered_lines": file_data.get("uncovered_lines", []),
                    "impact": impact,
                    "uncovered_statements": uncovered_statements
                })
        
        # Sort files by impact (highest first)
        uncovered_files.sort(key=lambda x: x["impact"], reverse=True)
        
        logger.info(f"Prioritized {len(uncovered_files)} files for testing in {repo_path}")
        return uncovered_files
    
    def _generate_test_file_path(self, file_path: str, language: str) -> str:
        """
        Generate the test file path for a source file.
        
        Args:
            file_path: Path to the source file
            language: Programming language
            
        Returns:
            str: Test file path
        """
        dir_path = os.path.dirname(file_path)
        file_name = os.path.basename(file_path)
        base_name, ext = os.path.splitext(file_name)
        
        if language == "go":
            return os.path.join(dir_path, f"{base_name}_test.go")
        elif language == "python":
            return os.path.join(dir_path, f"test_{base_name}.py")
        elif language == "typescript":
            return os.path.join(dir_path, f"{base_name}.test.ts")
        elif language == "java":
            return os.path.join(dir_path, f"{base_name}Test.java")
        elif language == "rust":
            return os.path.join(dir_path, f"{base_name}_test.rs")
        elif language == "cpp":
            return os.path.join(dir_path, f"{base_name}_test.cpp")
        elif language == "csharp":
            return os.path.join(dir_path, f"{base_name}Test.cs")
        elif language == "php":
            return os.path.join(dir_path, f"{base_name}Test.php")
        else:
            return os.path.join(dir_path, f"{base_name}_test{ext}")
    
    def _extract_functions(self, file_content: str, language: str, 
                          function_name: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Extract functions from a file.
        
        Args:
            file_content: Content of the file
            language: Programming language
            function_name: Name of the function to extract
            
        Returns:
            List[Dict[str, Any]]: Extracted functions
        """
        # This is a simplified implementation
        # In a real implementation, we would use language-specific parsers
        
        functions = []
        
        # Simple regex-based extraction for demonstration
        import re
        
        if language == "go":
            # Match Go functions
            pattern = r"func\s+(\w+)\s*\((.*?)\)\s*(?:\(?(.*?)\)?)?\s*\{"
            matches = re.finditer(pattern, file_content, re.DOTALL)
            
            for match in matches:
                name = match.group(1)
                params = match.group(2)
                returns = match.group(3) or ""
                
                if function_name and name != function_name:
                    continue
                
                functions.append({
                    "name": name,
                    "params": params,
                    "returns": returns,
                    "language": language
                })
        
        elif language == "python":
            # Match Python functions
            pattern = r"def\s+(\w+)\s*\((.*?)\)\s*(?:->)?\s*(.*?):"
            matches = re.finditer(pattern, file_content, re.DOTALL)
            
            for match in matches:
                name = match.group(1)
                params = match.group(2)
                returns = match.group(3) or ""
                
                if function_name and name != function_name:
                    continue
                
                functions.append({
                    "name": name,
                    "params": params,
                    "returns": returns,
                    "language": language
                })
        
        # Add more language-specific extractors as needed
        
        return functions
    
    def _generate_test_case(self, function: Dict[str, Any], language: str, 
                           test_framework: str) -> str:
        """
        Generate a test case for a function.
        
        Args:
            function: Function information
            language: Programming language
            test_framework: Test framework
            
        Returns:
            str: Generated test case
        """
        name = function["name"]
        
        if language == "go":
            return f"""func Test{name.capitalize()}(t *testing.T) {{
    // TODO: Implement test for {name}
    // Example:
    // input := ...
    // expected := ...
    // result := {name}(input)
    // if result != expected {{
    //     t.Errorf("Expected %v, got %v", expected, result)
    // }}
}}"""
        
        elif language == "python":
            if test_framework == "pytest":
                return f"""def test_{name}():
    # TODO: Implement test for {name}
    # Example:
    # input_value = ...
    # expected = ...
    # result = {name}(input_value)
    # assert result == expected"""
            else:
                return f"""def test_{name}(self):
    # TODO: Implement test for {name}
    # Example:
    # input_value = ...
    # expected = ...
    # result = {name}(input_value)
    # self.assertEqual(result, expected)"""
        
        elif language == "typescript":
            return f"""test('{name} should work correctly', () => {{
    // TODO: Implement test for {name}
    // Example:
    // const input = ...;
    // const expected = ...;
    // const result = {name}(input);
    // expect(result).toEqual(expected);
}});"""
        
        # Default case for other languages
        return f"// TODO: Implement test for {name}"
    
    def _fill_template(self, template: str, data: Dict[str, Any]) -> str:
        """
        Fill in a template with data.
        
        Args:
            template: Template string
            data: Data to fill in
            
        Returns:
            str: Filled template
        """
        result = template
        
        for key, value in data.items():
            result = result.replace(f"{{{{ {key} }}}}", str(value))
        
        return result
    
    def _get_default_template(self, language: str, test_framework: str) -> str:
        """
        Get a default test template for a language and framework.
        
        Args:
            language: Programming language
            test_framework: Test framework
            
        Returns:
            str: Default test template
        """
        if language == "go":
            return """package {{ import_path }}_test

import (
    "testing"
    
    // Import the package being tested
    "{{ import_path }}"
)

{{ test_cases }}
"""
        
        elif language == "python":
            if test_framework == "pytest":
                return """# Test file for {{ file_name }}

import pytest
# Import the module being tested
from {{ import_path }} import *

{{ test_cases }}
"""
            else:
                return """# Test file for {{ file_name }}

import unittest
# Import the module being tested
from {{ import_path }} import *

class Test{{ import_path.capitalize() }}(unittest.TestCase):
    {{ test_cases }}

if __name__ == '__main__':
    unittest.main()
"""
        
        elif language == "typescript":
            return """// Test file for {{ file_name }}

import { {{ import_path }} } from './{{ import_path }}';

{{ test_cases }}
"""
        
        # Add more language-specific templates as needed
        
        return "// TODO: Implement tests for {{ file_name }}"
    
    def _is_test_file(self, file_path: str, language: str) -> bool:
        """
        Check if a file is a test file.
        
        Args:
            file_path: Path to the file
            language: Programming language
            
        Returns:
            bool: True if the file is a test file, False otherwise
        """
        file_name = os.path.basename(file_path)
        
        if language == "go":
            return file_name.endswith("_test.go")
        elif language == "python":
            return file_name.startswith("test_") or file_name.endswith("_test.py")
        elif language == "typescript":
            return ".test." in file_name or ".spec." in file_name
        elif language == "java":
            return file_name.endswith("Test.java")
        elif language == "rust":
            return file_name.endswith("_test.rs")
        elif language == "cpp":
            return file_name.endswith("_test.cpp") or file_name.endswith("Test.cpp")
        elif language == "csharp":
            return file_name.endswith("Test.cs")
        elif language == "php":
            return file_name.endswith("Test.php")
        
        return False
