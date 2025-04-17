"""
Code Execution Module for Teemo

This module provides integration with the LibreChat Code Interpreter API,
enabling code execution capabilities across multiple programming languages.
"""

import os
import json
from typing import Dict, List, Optional, Any, Union

from .librechat_api_client import LibreChatCodeInterpreter

class CodeExecutionModule:
    """Module for executing code across multiple programming languages."""
    
    def __init__(self):
        """Initialize the Code Execution Module."""
        try:
            api_key = os.environ.get("LIBRECHAT_CODE_API_KEY")
            self.client = LibreChatCodeInterpreter(api_key)
            self.supported_languages = self.client.get_supported_languages()
        except Exception as e:
            print(f"Error initializing LibreChat Code Interpreter: {str(e)}")
            self.client = None
            self.supported_languages = [
                "python", "javascript", "typescript", 
                "c++", "c#", "go", "rust", "php"
            ]
    
    def execute_code(self, code: str, language: str = "python", timeout: int = 30) -> Dict[str, Any]:
        """Execute code using the LibreChat Code Interpreter.
        
        Args:
            code: The code to execute
            language: Programming language
            timeout: Maximum execution time in seconds
            
        Returns:
            Dict containing execution results
        """
        if not self.client:
            return {
                "success": False,
                "error": "LibreChat Code Interpreter client not initialized",
                "output": "",
                "execution_time": 0
            }
        
        if language not in self.supported_languages:
            return {
                "success": False,
                "error": f"Unsupported language: {language}. Supported languages: {', '.join(self.supported_languages)}",
                "output": "",
                "execution_time": 0
            }
        
        return self.client.execute_code(code, language, timeout)
    
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
        if not self.client:
            return {
                "success": False,
                "error": "LibreChat Code Interpreter client not initialized",
                "output": "",
                "execution_time": 0
            }
        
        if language not in self.supported_languages:
            return {
                "success": False,
                "error": f"Unsupported language: {language}. Supported languages: {', '.join(self.supported_languages)}",
                "output": "",
                "execution_time": 0
            }
        
        return self.client.execute_with_dependencies(code, language, dependencies, timeout)
    
    def get_supported_languages(self) -> List[str]:
        """Get list of supported programming languages.
        
        Returns:
            List of supported language strings
        """
        return self.supported_languages
