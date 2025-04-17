"""
Base Module for Corki Agent - Backup Agent for 100% Code Coverage

This module defines the BaseModule abstract class that all Corki modules must inherit from.
It provides common functionality and interface requirements for all modules.
"""

import abc
import json
import logging
import os
from typing import Any, Dict, List, Optional, Tuple, Union

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class BaseModule(abc.ABC):
    """
    Abstract base class for all Corki agent modules.
    
    All modules in the Corki agent must inherit from this class and implement
    its abstract methods. This ensures a consistent interface across all modules.
    """
    
    def __init__(self, name: str, config: Optional[Dict[str, Any]] = None):
        """
        Initialize the base module.
        
        Args:
            name: The name of the module
            config: Optional configuration dictionary for the module
        """
        self.name = name
        self.config = config or {}
        self.tools = {}
        self.initialized = False
        logger.info(f"Initializing {name} module")
    
    @abc.abstractmethod
    def initialize(self) -> bool:
        """
        Initialize the module and register its tools.
        
        This method must be implemented by all subclasses. It should set up
        any resources needed by the module and register all tools provided
        by the module.
        
        Returns:
            bool: True if initialization was successful, False otherwise
        """
        pass
    
    @abc.abstractmethod
    def cleanup(self) -> bool:
        """
        Clean up any resources used by the module.
        
        This method must be implemented by all subclasses. It should release
        any resources acquired during initialization or operation.
        
        Returns:
            bool: True if cleanup was successful, False otherwise
        """
        pass
    
    def register_tool(self, tool_name: str, tool_function: Any, tool_description: str, 
                     tool_args: List[Dict[str, Any]], tool_returns: Dict[str, Any]) -> bool:
        """
        Register a tool provided by this module.
        
        Args:
            tool_name: The name of the tool
            tool_function: The function that implements the tool
            tool_description: A description of what the tool does
            tool_args: A list of dictionaries describing the tool's arguments
            tool_returns: A dictionary describing the tool's return value
            
        Returns:
            bool: True if the tool was registered successfully, False otherwise
        """
        if tool_name in self.tools:
            logger.warning(f"Tool {tool_name} already registered, overwriting")
        
        self.tools[tool_name] = {
            "function": tool_function,
            "description": tool_description,
            "args": tool_args,
            "returns": tool_returns
        }
        
        logger.info(f"Registered tool {tool_name} in module {self.name}")
        return True
    
    def get_tool(self, tool_name: str) -> Optional[Dict[str, Any]]:
        """
        Get a tool by name.
        
        Args:
            tool_name: The name of the tool to get
            
        Returns:
            Optional[Dict[str, Any]]: The tool if found, None otherwise
        """
        return self.tools.get(tool_name)
    
    def get_all_tools(self) -> Dict[str, Dict[str, Any]]:
        """
        Get all tools registered by this module.
        
        Returns:
            Dict[str, Dict[str, Any]]: A dictionary of all tools
        """
        return self.tools
    
    def execute_tool(self, tool_name: str, *args, **kwargs) -> Tuple[bool, Any]:
        """
        Execute a tool by name.
        
        Args:
            tool_name: The name of the tool to execute
            *args: Positional arguments to pass to the tool
            **kwargs: Keyword arguments to pass to the tool
            
        Returns:
            Tuple[bool, Any]: A tuple containing a success flag and the result
                             of the tool execution if successful, or an error
                             message if not
        """
        tool = self.get_tool(tool_name)
        if not tool:
            return False, f"Tool {tool_name} not found in module {self.name}"
        
        try:
            result = tool["function"](*args, **kwargs)
            return True, result
        except Exception as e:
            logger.error(f"Error executing tool {tool_name}: {str(e)}")
            return False, str(e)
    
    def load_config_from_file(self, config_path: str) -> bool:
        """
        Load configuration from a JSON file.
        
        Args:
            config_path: Path to the configuration file
            
        Returns:
            bool: True if the configuration was loaded successfully, False otherwise
        """
        try:
            if not os.path.exists(config_path):
                logger.error(f"Configuration file {config_path} does not exist")
                return False
            
            with open(config_path, 'r') as f:
                self.config = json.load(f)
            
            logger.info(f"Loaded configuration from {config_path}")
            return True
        except Exception as e:
            logger.error(f"Error loading configuration from {config_path}: {str(e)}")
            return False
    
    def save_config_to_file(self, config_path: str) -> bool:
        """
        Save configuration to a JSON file.
        
        Args:
            config_path: Path to the configuration file
            
        Returns:
            bool: True if the configuration was saved successfully, False otherwise
        """
        try:
            with open(config_path, 'w') as f:
                json.dump(self.config, f, indent=2)
            
            logger.info(f"Saved configuration to {config_path}")
            return True
        except Exception as e:
            logger.error(f"Error saving configuration to {config_path}: {str(e)}")
            return False
