"""
LibreChat Code Interpreter API Client for Teemo

This module provides integration with the LibreChat Code Interpreter API,
enabling code execution capabilities across multiple programming languages.
"""

import os
import requests
import json
from typing import Dict, List, Optional, Union, Any

class LibreChatCodeInterpreter:
    """Client for the LibreChat Code Interpreter API."""
    
    def __init__(self, api_key: Optional[str] = None):
        """Initialize the LibreChat Code Interpreter client.
        
        Args:
            api_key: API key for authentication. If not provided, will attempt to read from environment.
        """
        self.api_key = api_key or os.environ.get("LIBRECHAT_CODE_API_KEY")
        if not self.api_key:
            raise ValueError("LibreChat API key is required")
        
        self.base_url = "http://185.192.220.224:8000/api/v1"
        self.headers = {
            "Content-Type": "application/json",
            "x-api-key": self.api_key
        }
    
    def execute_code(self, code: str, language: str = "python", timeout: int = 30) -> Dict[str, Any]:
        """Execute code using the LibreChat Code Interpreter.
        
        Args:
            code: The code to execute
            language: Programming language (python, javascript, typescript, c++, c#, go, rust, php)
            timeout: Maximum execution time in seconds
            
        Returns:
            Dict containing execution results with keys:
            - success: Boolean indicating if execution was successful
            - output: String containing execution output
            - error: Error message if execution failed
            - execution_time: Time taken to execute in seconds
        """
        endpoint = f"{self.base_url}/execute"
        payload = {
            "code": code,
            "language": language,
            "timeout": timeout
        }
        
        try:
            response = requests.post(endpoint, headers=self.headers, json=payload)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            return {
                "success": False,
                "error": f"API request failed: {str(e)}",
                "output": "",
                "execution_time": 0
            }
    
    def get_supported_languages(self) -> List[str]:
        """Get list of supported programming languages.
        
        Returns:
            List of supported language strings
        """
        endpoint = f"{self.base_url}/languages"
        
        try:
            response = requests.get(endpoint, headers=self.headers)
            response.raise_for_status()
            return response.json().get("languages", [])
        except requests.exceptions.RequestException:
            return ["python", "javascript", "typescript", "c++", "c#", "go", "rust", "php"]
    
    def execute_with_dependencies(self, code: str, language: str = "python", 
                                 dependencies: List[str] = None, 
                                 timeout: int = 60) -> Dict[str, Any]:
        """Execute code with specified dependencies.
        
        Args:
            code: The code to execute
            language: Programming language
            dependencies: List of dependencies to install before execution
            timeout: Maximum execution time in seconds
            
        Returns:
            Dict containing execution results
        """
        endpoint = f"{self.base_url}/execute_with_deps"
        payload = {
            "code": code,
            "language": language,
            "dependencies": dependencies or [],
            "timeout": timeout
        }
        
        try:
            response = requests.post(endpoint, headers=self.headers, json=payload)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            return {
                "success": False,
                "error": f"API request failed: {str(e)}",
                "output": "",
                "execution_time": 0
            }
