"""
Plugin manager for Neovim in agent_runtime.

This script manages Neovim plugins for the agent_runtime system.
"""

import os
import sys
import subprocess
import json
import yaml
import logging
from typing import Dict, Any, List, Optional

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger("neovim-plugin-manager")


class NeovimPluginManager:
    """Manager for Neovim plugins in agent_runtime."""

    def __init__(self, config_path: str, plugins_dir: str):
        """Initialize the plugin manager.
        
        Args:
            config_path: Path to the config.yaml file
            plugins_dir: Path to the plugins directory
        """
        self.config_path = config_path
        self.plugins_dir = plugins_dir
        self.config = self._load_config()
        self.neovim_config_dir = os.path.expanduser("~/.config/nvim")
        self.plugin_configs = {}

    def _load_config(self) -> Dict[str, Any]:
        """Load the config.yaml file.
        
        Returns:
            Dict[str, Any]: Configuration data
        """
        try:
            with open(self.config_path, "r") as f:
                return yaml.safe_load(f)
        except Exception as e:
            logger.error(f"Error loading config: {e}")
            return {}

    def get_plugins(self) -> List[Dict[str, Any]]:
        """Get the list of plugins from the config.
        
        Returns:
            List[Dict[str, Any]]: List of plugin configurations
        """
        return self.config.get("plugins", [])

    def install_plugin(self, plugin_name: str) -> bool:
        """Install a specific plugin.
        
        Args:
            plugin_name: Name of the plugin to install
            
        Returns:
            bool: True if installation was successful
        """
        plugins = self.get_plugins()
        plugin = next((p for p in plugins if p["name"] == plugin_name), None)
        
        if not plugin:
            logger.error(f"Plugin {plugin_name} not found in config")
            return False
        
        logger.info(f"Installing plugin: {plugin_name} from {plugin['url']}")
        
        self.plugin_configs[plugin_name] = plugin
        
        return True

    def install_all_plugins(self) -> bool:
        """Install all plugins from the config.
        
        Returns:
            bool: True if all installations were successful
        """
        plugins = self.get_plugins()
        success = True
        
        for plugin in plugins:
            plugin_name = plugin["name"]
            if not self.install_plugin(plugin_name):
                success = False
        
        return success

    def run_install_script(self) -> bool:
        """Run the install.sh script.
        
        Returns:
            bool: True if the script ran successfully
        """
        install_script = os.path.join(self.plugins_dir, "install.sh")
        
        if not os.path.exists(install_script):
            logger.error(f"Install script not found: {install_script}")
            return False
        
        try:
            logger.info(f"Running install script: {install_script}")
            subprocess.run(["bash", install_script], check=True)
            return True
        except subprocess.CalledProcessError as e:
            logger.error(f"Error running install script: {e}")
            return False

    def setup_neovim_config(self) -> bool:
        """Set up the Neovim configuration.
        
        Returns:
            bool: True if setup was successful
        """
        try:
            os.makedirs(self.neovim_config_dir, exist_ok=True)
            
            init_lua = os.path.join(self.plugins_dir, "init.lua")
            if os.path.exists(init_lua):
                logger.info(f"Copying {init_lua} to {self.neovim_config_dir}")
                with open(init_lua, "r") as src, open(os.path.join(self.neovim_config_dir, "init.lua"), "w") as dst:
                    dst.write(src.read())
            
            return True
        except Exception as e:
            logger.error(f"Error setting up Neovim config: {e}")
            return False

    def verify_installation(self) -> Dict[str, bool]:
        """Verify that all plugins are installed correctly.
        
        Returns:
            Dict[str, bool]: Status of each plugin
        """
        plugins = self.get_plugins()
        status = {}
        
        for plugin in plugins:
            plugin_name = plugin["name"]
            status[plugin_name] = plugin_name in self.plugin_configs
        
        return status

    def get_plugin_commands(self, plugin_name: str) -> List[Dict[str, str]]:
        """Get the commands for a specific plugin.
        
        Args:
            plugin_name: Name of the plugin
            
        Returns:
            List[Dict[str, str]]: List of commands
        """
        plugins = self.get_plugins()
        plugin = next((p for p in plugins if p["name"] == plugin_name), None)
        
        if not plugin:
            logger.error(f"Plugin {plugin_name} not found in config")
            return []
        
        return plugin.get("commands", [])

    def get_plugin_info(self, plugin_name: str) -> Optional[Dict[str, Any]]:
        """Get information about a specific plugin.
        
        Args:
            plugin_name: Name of the plugin
            
        Returns:
            Optional[Dict[str, Any]]: Plugin information
        """
        plugins = self.get_plugins()
        return next((p for p in plugins if p["name"] == plugin_name), None)


def main():
    """Main entry point for the plugin manager."""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    neovim_dir = os.path.dirname(script_dir)
    
    config_path = os.path.join(neovim_dir, "config.yaml")
    
    manager = NeovimPluginManager(config_path, script_dir)
    
    if len(sys.argv) > 1:
        command = sys.argv[1]
        
        if command == "install":
            if len(sys.argv) > 2:
                plugin_name = sys.argv[2]
                success = manager.install_plugin(plugin_name)
                print(f"Installation of {plugin_name}: {'Successful' if success else 'Failed'}")
            else:
                success = manager.install_all_plugins()
                print(f"Installation of all plugins: {'Successful' if success else 'Failed'}")
        
        elif command == "setup":
            success = manager.setup_neovim_config()
            print(f"Setup of Neovim config: {'Successful' if success else 'Failed'}")
        
        elif command == "run-install-script":
            success = manager.run_install_script()
            print(f"Running install script: {'Successful' if success else 'Failed'}")
        
        elif command == "verify":
            status = manager.verify_installation()
            print("Plugin installation status:")
            for plugin, installed in status.items():
                print(f"  {plugin}: {'Installed' if installed else 'Not installed'}")
        
        elif command == "list":
            plugins = manager.get_plugins()
            print("Available plugins:")
            for plugin in plugins:
                print(f"  {plugin['name']}: {plugin['description']}")
                print(f"    URL: {plugin['url']}")
                print("    Commands:")
                for cmd in plugin.get("commands", []):
                    print(f"      {cmd['name']}: {cmd['description']}")
                    print(f"        Usage: {cmd['usage']}")
                print()
        
        elif command == "info":
            if len(sys.argv) > 2:
                plugin_name = sys.argv[2]
                plugin = manager.get_plugin_info(plugin_name)
                if plugin:
                    print(f"Plugin: {plugin['name']}")
                    print(f"Description: {plugin['description']}")
                    print(f"URL: {plugin['url']}")
                    print("Commands:")
                    for cmd in plugin.get("commands", []):
                        print(f"  {cmd['name']}: {cmd['description']}")
                        print(f"    Usage: {cmd['usage']}")
                else:
                    print(f"Plugin {plugin_name} not found")
            else:
                print("Please specify a plugin name")
        
        else:
            print(f"Unknown command: {command}")
            print("Available commands: install, setup, run-install-script, verify, list, info")
    else:
        print("Please specify a command")
        print("Available commands: install, setup, run-install-script, verify, list, info")


if __name__ == "__main__":
    main()
