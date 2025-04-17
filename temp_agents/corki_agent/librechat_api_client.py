"""
LibreChat API Client for Corki Agent - Backup Agent for 100% Code Coverage

This module provides a client for interacting with the LibreChat Code Interpreter API.
It enables the agent to execute code in multiple languages and analyze results.
"""

import json
import logging
import os
import requests
from typing import Any, Dict, List, Optional, Union

from base_module import BaseModule

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class LibreChatAPIClient(BaseModule):
    """
    Client for interacting with the LibreChat Code Interpreter API.
    
    This module provides tools for executing code in multiple languages
    and analyzing the results using the LibreChat Code Interpreter API.
    """
    
    def __init__(self, config: Optional[Dict[str, Any]] = None):
        """
        Initialize the LibreChat API client.
        
        Args:
            config: Optional configuration dictionary for the module
        """
        super().__init__("librechat_api_client", config)
        self.api_key = os.environ.get("LIBRECHAT_CODE_API_KEY", "")
        self.api_url = self.config.get("api_url", "https://librechat.ai/api/code")
        self.supported_languages = [
            "python", "go", "typescript", "javascript", "java", 
            "rust", "cpp", "csharp", "php", "bash", "sql"
        ]
    
    def initialize(self) -> bool:
        """
        Initialize the LibreChat API client and register its tools.
        
        Returns:
            bool: True if initialization was successful, False otherwise
        """
        if not self.api_key:
            logger.error("LibreChat API key not found in environment variables")
            return False
        
        self.register_tool(
            "execute_code",
            self.execute_code,
            "Execute code using the LibreChat Code Interpreter API",
            [
                {
                    "name": "code",
                    "type": "string",
                    "description": "Code to execute"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
                },
                {
                    "name": "timeout",
                    "type": "number",
                    "description": "Timeout in seconds",
                    "required": False
                }
            ],
            {
                "type": "object",
                "description": "Execution results"
            }
        )
        
        self.register_tool(
            "generate_test",
            self.generate_test,
            "Generate a test for a given function using the LibreChat Code Interpreter API",
            [
                {
                    "name": "function_code",
                    "type": "string",
                    "description": "Function code to test"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
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
            "analyze_code",
            self.analyze_code,
            "Analyze code using the LibreChat Code Interpreter API",
            [
                {
                    "name": "code",
                    "type": "string",
                    "description": "Code to analyze"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
                },
                {
                    "name": "analysis_type",
                    "type": "string",
                    "description": "Type of analysis to perform",
                    "required": False
                }
            ],
            {
                "type": "object",
                "description": "Analysis results"
            }
        )
        
        self.register_tool(
            "fix_code",
            self.fix_code,
            "Fix code issues using the LibreChat Code Interpreter API",
            [
                {
                    "name": "code",
                    "type": "string",
                    "description": "Code to fix"
                },
                {
                    "name": "language",
                    "type": "string",
                    "description": "Programming language"
                },
                {
                    "name": "issues",
                    "type": "array",
                    "description": "List of issues to fix",
                    "required": False
                }
            ],
            {
                "type": "object",
                "description": "Fixed code"
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
    
    def execute_code(self, code: str, language: str, timeout: Optional[int] = None) -> Dict[str, Any]:
        """
        Execute code using the LibreChat Code Interpreter API.
        
        Args:
            code: Code to execute
            language: Programming language
            timeout: Timeout in seconds
            
        Returns:
            Dict[str, Any]: Execution results
        """
        if not self.initialized:
            logger.error("LibreChat API client not initialized")
            return {"error": "LibreChat API client not initialized"}
        
        if language not in self.supported_languages:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        try:
            data = {
                "code": code,
                "language": language
            }
            
            if timeout:
                data["timeout"] = str(timeout)
            
            headers = {
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.api_key}"
            }
            
            response = requests.post(
                f"{self.api_url}/execute",
                headers=headers,
                json=data
            )
            
            if response.status_code != 200:
                logger.error(f"API request failed: {response.status_code} {response.text}")
                return {"error": f"API request failed: {response.status_code} {response.text}"}
            
            result = response.json()
            
            logger.info(f"Executed code in {language}")
            return result
        
        except Exception as e:
            logger.error(f"Error executing code: {str(e)}")
            return {"error": f"Error executing code: {str(e)}"}
    
    def generate_test(self, function_code: str, language: str, 
                     test_framework: Optional[str] = None) -> Dict[str, Any]:
        """
        Generate a test for a given function using the LibreChat Code Interpreter API.
        
        Args:
            function_code: Function code to test
            language: Programming language
            test_framework: Test framework to use
            
        Returns:
            Dict[str, Any]: Generated test
        """
        if not self.initialized:
            logger.error("LibreChat API client not initialized")
            return {"error": "LibreChat API client not initialized"}
        
        if language not in self.supported_languages:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        try:
            data = {
                "function_code": function_code,
                "language": language
            }
            
            if test_framework:
                data["test_framework"] = test_framework
            
            headers = {
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.api_key}"
            }
            
            response = requests.post(
                f"{self.api_url}/generate_test",
                headers=headers,
                json=data
            )
            
            if response.status_code != 200:
                logger.error(f"API request failed: {response.status_code} {response.text}")
                return {"error": f"API request failed: {response.status_code} {response.text}"}
            
            result = response.json()
            
            logger.info(f"Generated test for {language} function")
            return result
        
        except Exception as e:
            logger.error(f"Error generating test: {str(e)}")
            return {"error": f"Error generating test: {str(e)}"}
    
    def analyze_code(self, code: str, language: str, 
                    analysis_type: Optional[str] = None) -> Dict[str, Any]:
        """
        Analyze code using the LibreChat Code Interpreter API.
        
        Args:
            code: Code to analyze
            language: Programming language
            analysis_type: Type of analysis to perform
            
        Returns:
            Dict[str, Any]: Analysis results
        """
        if not self.initialized:
            logger.error("LibreChat API client not initialized")
            return {"error": "LibreChat API client not initialized"}
        
        if language not in self.supported_languages:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        try:
            data = {
                "code": code,
                "language": language
            }
            
            if analysis_type:
                data["analysis_type"] = analysis_type
            
            headers = {
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.api_key}"
            }
            
            response = requests.post(
                f"{self.api_url}/analyze",
                headers=headers,
                json=data
            )
            
            if response.status_code != 200:
                logger.error(f"API request failed: {response.status_code} {response.text}")
                return {"error": f"API request failed: {response.status_code} {response.text}"}
            
            result = response.json()
            
            logger.info(f"Analyzed code in {language}")
            return result
        
        except Exception as e:
            logger.error(f"Error analyzing code: {str(e)}")
            return {"error": f"Error analyzing code: {str(e)}"}
    
    def fix_code(self, code: str, language: str, 
                issues: Optional[List[Dict[str, Any]]] = None) -> Dict[str, Any]:
        """
        Fix code issues using the LibreChat Code Interpreter API.
        
        Args:
            code: Code to fix
            language: Programming language
            issues: List of issues to fix
            
        Returns:
            Dict[str, Any]: Fixed code
        """
        if not self.initialized:
            logger.error("LibreChat API client not initialized")
            return {"error": "LibreChat API client not initialized"}
        
        if language not in self.supported_languages:
            logger.error(f"Unsupported language: {language}")
            return {"error": f"Unsupported language: {language}"}
        
        try:
            data = {
                "code": code,
                "language": language
            }
            
            if issues:
                data["issues"] = json.dumps(issues)
            
            headers = {
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.api_key}"
            }
            
            response = requests.post(
                f"{self.api_url}/fix",
                headers=headers,
                json=data
            )
            
            if response.status_code != 200:
                logger.error(f"API request failed: {response.status_code} {response.text}")
                return {"error": f"API request failed: {response.status_code} {response.text}"}
            
            result = response.json()
            
            logger.info(f"Fixed code in {language}")
            return result
        
        except Exception as e:
            logger.error(f"Error fixing code: {str(e)}")
            return {"error": f"Error fixing code: {str(e)}"}
