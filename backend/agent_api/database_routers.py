"""
Database routers for the agent_api project.

This module provides database routers for directing database operations
to the appropriate database based on the model.
"""

import logging

logger = logging.getLogger(__name__)


class AgentRuntimeRouter:
    """
    Main router for the agent_runtime project.
    
    This router directs database operations to the appropriate database
    based on the model's metadata.
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
            return 'agent_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            return 'trajectory_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            return 'ml_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'analytics_model'):
            return 'default'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'supabase_db'):
            db_name = getattr(model._meta, 'supabase_db')
            if db_name:
                return db_name
        
        if model._meta.app_label in ['api', 'ml_api', 'python_agent', 'python_ml', 'app']:
            return 'default'
        
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
            return 'agent_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            return 'trajectory_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            return 'ml_db'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'analytics_model'):
            return 'default'
        
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'supabase_db'):
            db_name = getattr(model._meta, 'supabase_db')
            if db_name:
                return db_name
        
        if model._meta.app_label in ['api', 'ml_api', 'python_agent', 'python_ml', 'app']:
            return 'default'
        
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
        if hasattr(obj1._meta, 'agent_model') and hasattr(obj2._meta, 'agent_model'):
            return True
        
        if hasattr(obj1._meta, 'trajectory_model') and hasattr(obj2._meta, 'trajectory_model'):
            return True
        
        if hasattr(obj1._meta, 'ml_model') and hasattr(obj2._meta, 'ml_model'):
            return True
        
        if hasattr(obj1._meta, 'analytics_model') and hasattr(obj2._meta, 'analytics_model'):
            return True
        
        if (hasattr(obj1._meta, 'supabase_db') and hasattr(obj2._meta, 'supabase_db') and
                getattr(obj1._meta, 'supabase_db') == getattr(obj2._meta, 'supabase_db')):
            return True
        
        if obj1._meta.app_label == obj2._meta.app_label:
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
        model = hints.get('model')
        
        if db == 'agent_db' and app_label == 'python_agent':
            if model and hasattr(model._meta, 'agent_model'):
                return True
            elif model_name and model_name.lower().startswith('agent'):
                return True
            return False
        
        if db == 'trajectory_db' and app_label == 'python_agent':
            if model and hasattr(model._meta, 'trajectory_model'):
                return True
            elif model_name and model_name.lower().startswith('trajectory'):
                return True
            return False
        
        if db == 'ml_db' and app_label == 'python_agent':
            if model and hasattr(model._meta, 'ml_model'):
                return True
            elif model_name and model_name.lower().startswith('ml'):
                return True
            return False
        
        if db == 'default':
            if model and hasattr(model._meta, 'analytics_model'):
                return True
            elif app_label in ['api', 'ml_api', 'app']:
                return True
            elif app_label == 'python_agent' and model_name and model_name.lower().startswith('analytics'):
                return True
        
        if db.startswith('supabase_'):
            if model and hasattr(model._meta, 'supabase_db') and getattr(model._meta, 'supabase_db') == db:
                return True
            return False
        
        return None


class AgentRouter:
    """Legacy router for agent-related models."""
    
    def db_for_read(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'agent_model'):
            logger.warning("Using legacy AgentRouter, please update to AgentRuntimeRouter")
            return 'agent_db'
        return None
    
    def db_for_write(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'agent_model'):
            logger.warning("Using legacy AgentRouter, please update to AgentRuntimeRouter")
            return 'agent_db'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'agent_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'agent_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        if db == 'agent_db' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'agent_model'):
                return True
        return None


class TrajectoryRouter:
    """Legacy router for trajectory-related models."""
    
    def db_for_read(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            logger.warning("Using legacy TrajectoryRouter, please update to AgentRuntimeRouter")
            return 'trajectory_db'
        return None
    
    def db_for_write(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'trajectory_model'):
            logger.warning("Using legacy TrajectoryRouter, please update to AgentRuntimeRouter")
            return 'trajectory_db'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'trajectory_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'trajectory_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        if db == 'trajectory_db' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'trajectory_model'):
                return True
        return None


class MLRouter:
    """Legacy router for ML-related models."""
    
    def db_for_read(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            logger.warning("Using legacy MLRouter, please update to AgentRuntimeRouter")
            return 'ml_db'
        return None
    
    def db_for_write(self, model, **hints):
        if model._meta.app_label == 'python_agent' and hasattr(model._meta, 'ml_model'):
            logger.warning("Using legacy MLRouter, please update to AgentRuntimeRouter")
            return 'ml_db'
        return None
    
    def allow_relation(self, obj1, obj2, **hints):
        if (
            obj1._meta.app_label == 'python_agent' and hasattr(obj1._meta, 'ml_model') and
            obj2._meta.app_label == 'python_agent' and hasattr(obj2._meta, 'ml_model')
        ):
            return True
        return None
    
    def allow_migrate(self, db, app_label, model_name=None, **hints):
        if db == 'ml_db' and app_label == 'python_agent':
            model = hints.get('model')
            if model and hasattr(model._meta, 'ml_model'):
                return True
        return None
