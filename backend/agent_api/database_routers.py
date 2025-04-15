"""
Database routers for the agent_api project.

This module provides database routers for directing database operations
to the appropriate database based on the model.
"""

class AgentRouter:
    """
    Router for agent-related models.
    
    This router directs database operations for agent-related models
    to the agent database.
    """
    
    def db_for_read(self, model, **hints):
        """
        Suggest the database that should be used for read operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'agent_model'):
            return 'agent'
        return None
    
    def db_for_write(self, model, **hints):
        """
        Suggest the database that should be used for write operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'agent_model'):
            return 'agent'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        """
        Determine if a relation should be allowed between two objects.
        
        Args:
            obj1: The first object
            obj2: The second object
            hints: Additional hints
            
        Returns:
            bool: True if the relation should be allowed, False otherwise
        """
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'agent_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'agent_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        """
        Determine if a migration should be run on a database.
        
        Args:
            db: The database to check
            app_label: The app label
            model_name: The model name
            hints: Additional hints
            
        Returns:
            bool: True if the migration should be run, False otherwise
        """
        if db == 'agent' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'agent_model'):
                return True
        return None


class TrajectoryRouter:
    """
    Router for trajectory-related models.
    
    This router directs database operations for trajectory-related models
    to the trajectory database.
    """
    
    def db_for_read(self, model, **hints):
        """
        Suggest the database that should be used for read operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            return 'trajectory'
        return None
    
    def db_for_write(self, model, **hints):
        """
        Suggest the database that should be used for write operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            return 'trajectory'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        """
        Determine if a relation should be allowed between two objects.
        
        Args:
            obj1: The first object
            obj2: The second object
            hints: Additional hints
            
        Returns:
            bool: True if the relation should be allowed, False otherwise
        """
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'trajectory_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'trajectory_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        """
        Determine if a migration should be run on a database.
        
        Args:
            db: The database to check
            app_label: The app label
            model_name: The model name
            hints: Additional hints
            
        Returns:
            bool: True if the migration should be run, False otherwise
        """
        if db == 'trajectory' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'trajectory_model'):
                return True
        return None


class MLRouter:
    """
    Router for ML-related models.
    
    This router directs database operations for ML-related models
    to the ML database.
    """
    
    def db_for_read(self, model, **hints):
        """
        Suggest the database that should be used for read operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            return 'ml'
        return None
    
    def db_for_write(self, model, **hints):
        """
        Suggest the database that should be used for write operations.
        
        Args:
            model: The model to suggest a database for
            hints: Additional hints
            
        Returns:
            str: The name of the database to use
        """
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            return 'ml'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        """
        Determine if a relation should be allowed between two objects.
        
        Args:
            obj1: The first object
            obj2: The second object
            hints: Additional hints
            
        Returns:
            bool: True if the relation should be allowed, False otherwise
        """
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'ml_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'ml_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        """
        Determine if a migration should be run on a database.
        
        Args:
            db: The database to check
            app_label: The app label
            model_name: The model name
            hints: Additional hints
            
        Returns:
            bool: True if the migration should be run, False otherwise
        """
        if db == 'ml' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'ml_model'):
                return True
        return None
