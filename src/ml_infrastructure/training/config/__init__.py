"""
Training Configuration

This package provides configuration classes for fine-tuning Llama 4 models.
"""

from .llama4_maverick_config import Llama4MaverickConfig, get_default_config as get_maverick_config
from .llama4_scout_config import Llama4ScoutConfig, get_default_config as get_scout_config

__all__ = [
    "Llama4MaverickConfig",
    "Llama4ScoutConfig",
    "get_maverick_config",
    "get_scout_config",
]
