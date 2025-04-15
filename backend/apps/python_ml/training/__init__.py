"""
Llama 4 Training

This package provides configuration, scripts, and utilities for fine-tuning
Llama 4 Maverick and Scout models on GitOps, Terraform, and Kubernetes issue data.
"""

from .config.llama4_maverick_config import (
    Llama4MaverickConfig,
    get_default_config as get_maverick_config,
)
from .config.llama4_scout_config import (
    Llama4ScoutConfig,
    get_default_config as get_scout_config,
)

__all__ = [
    "Llama4MaverickConfig",
    "Llama4ScoutConfig",
    "get_maverick_config",
    "get_scout_config",
]
