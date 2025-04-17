"""
Base Module for Teemo
Defines the abstract base class for all Teemo modules
"""

from abc import ABC, abstractmethod
from typing import Dict, List, Any, Optional

class BaseModule(ABC):
    """
    Abstract base class for all Teemo modules
    All modules must implement these abstract methods
    """
    
    @property
    @abstractmethod
    def name(self) -> str:
        """Returns the module name"""
        pass
    
    @property
    @abstractmethod
    def description(self) -> str:
        """Returns the module description"""
        pass
    
    @property
    @abstractmethod
    def tools(self) -> List[str]:
        """Returns a list of tools provided by the module"""
        pass
    
    @abstractmethod
    def initialize(self, context: Dict[str, Any]) -> bool:
        """
        Initializes the module with execution context
        
        Args:
            context: Dictionary containing execution context
            
        Returns:
            bool: True if initialization was successful, False otherwise
        """
        pass
    
    @abstractmethod
    def cleanup(self) -> bool:
        """
        Cleans up module resources
        
        Returns:
            bool: True if cleanup was successful, False otherwise
        """
        pass
