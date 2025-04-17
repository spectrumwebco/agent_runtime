"""
Neovim tool for the Python agent that allows interacting with Neovim as a
parallel workflow.
"""

import asyncio
import requests
from typing import Dict, Any, Optional, List

from ..utils.log import get_logger


class NeovimTool:
    """Tool for interacting with Neovim as a parallel workflow environment."""

    def __init__(self):
        """Initialize the Neovim tool."""
        self.logger = get_logger("neovim-tool", emoji="📝")
        self.neovim_api_base = "http://localhost:8090/neovim"
        self.active_instances = {}
        self.background_tasks = {}

    async def start_instance(self, instance_id: str) -> bool:
        """Start a new Neovim instance.

        Args:
            instance_id: Unique identifier for the Neovim instance

        Returns:
            bool: True if instance was started successfully
        """
        try:
            response = requests.post(
                f"{self.neovim_api_base}/start",
                json={"id": instance_id}
            )
            if response.status_code == 200:
                self.active_instances[instance_id] = True
                self.logger.info(f"Started Neovim instance: {instance_id}")
                return True
            else:
                self.logger.error(
                    f"Failed to start Neovim instance: {response.text}"
                )
                return False
        except Exception as e:
            self.logger.error(f"Error starting Neovim instance: {e}")
            return False

    async def stop_instance(self, instance_id: str) -> bool:
        """Stop a Neovim instance.

        Args:
            instance_id: Unique identifier for the Neovim instance

        Returns:
            bool: True if instance was stopped successfully
        """
        try:
            response = requests.post(
                f"{self.neovim_api_base}/stop",
                json={"id": instance_id}
            )
            if response.status_code == 200:
                if instance_id in self.active_instances:
                    del self.active_instances[instance_id]
                self.logger.info(f"Stopped Neovim instance: {instance_id}")
                return True
            else:
                self.logger.error(
                    f"Failed to stop Neovim instance: {response.text}"
                )
                return False
        except Exception as e:
            self.logger.error(f"Error stopping Neovim instance: {e}")
            return False

    async def execute_command(
        self,
        instance_id: str,
        command: str,
        file: Optional[str] = None,
        background: bool = False
    ) -> Optional[Dict[str, Any]]:
        """Execute a command in a Neovim instance.

        Args:
            instance_id: Unique identifier for the Neovim instance
            command: Neovim command to execute
            file: Optional file to open before executing the command
            background: Whether to run in background (parallel)

        Returns:
            Optional[Dict[str, Any]]: Command result or None if running in
            background
        """
        if instance_id not in self.active_instances:
            success = await self.start_instance(instance_id)
            if not success:
                return {"error": "Failed to start Neovim instance"}

        if file:
            file_open_result = await self._execute_command(
                instance_id, f":e {file}<CR>"
            )
            if "error" in file_open_result:
                return file_open_result

        if background:
            task = asyncio.create_task(
                self._execute_command(instance_id, command)
            )
            task_id = f"{instance_id}_{len(self.background_tasks)}"
            self.background_tasks[task_id] = task
            self.logger.info(
                f"Started background task {task_id} for Neovim instance: "
                f"{instance_id}"
            )
            return {"status": "running_in_background", "task_id": task_id}
        else:
            return await self._execute_command(instance_id, command)

    async def _execute_command(
        self, instance_id: str, command: str
    ) -> Dict[str, Any]:
        """Internal method to execute a command.

        Args:
            instance_id: Unique identifier for the Neovim instance
            command: Neovim command to execute

        Returns:
            Dict[str, Any]: Command result
        """
        try:
            response = requests.post(
                f"{self.neovim_api_base}/execute",
                json={"id": instance_id, "command": command}
            )
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(
                    f"Failed to execute Neovim command: {response.text}"
                )
                return {"error": f"Failed to execute command: {response.text}"}
        except Exception as e:
            self.logger.error(f"Error executing Neovim command: {e}")
            return {"error": f"Error executing command: {str(e)}"}

    async def get_state(self, instance_id: str) -> Dict[str, Any]:
        """Get the current state of a Neovim instance.

        Args:
            instance_id: Unique identifier for the Neovim instance

        Returns:
            Dict[str, Any]: Current state
        """
        try:
            response = requests.get(
                f"{self.neovim_api_base}/state",
                params={"id": instance_id}
            )
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(
                    f"Failed to get Neovim state: {response.text}"
                )
                return {"error": f"Failed to get state: {response.text}"}
        except Exception as e:
            self.logger.error(f"Error getting Neovim state: {e}")
            return {"error": f"Error getting state: {str(e)}"}

    async def get_background_tasks(self) -> List[str]:
        """Get a list of all background tasks.

        Returns:
            List[str]: List of task IDs
        """
        return list(self.background_tasks.keys())

    async def check_background_task(self, task_id: str) -> Dict[str, Any]:
        """Check the status of a background task.

        Args:
            task_id: ID of the background task

        Returns:
            Dict[str, Any]: Task status
        """
        if task_id not in self.background_tasks:
            return {"error": f"No background task with ID {task_id}"}

        task = self.background_tasks[task_id]
        if task.done():
            try:
                result = task.result()
                del self.background_tasks[task_id]
                return {"status": "completed", "result": result}
            except Exception as e:
                del self.background_tasks[task_id]
                return {"status": "failed", "error": str(e)}
        else:
            return {"status": "running"}

    async def cancel_background_task(self, task_id: str) -> Dict[str, Any]:
        """Cancel a background task.

        Args:
            task_id: ID of the background task

        Returns:
            Dict[str, Any]: Cancellation result
        """
        if task_id not in self.background_tasks:
            return {"error": f"No background task with ID {task_id}"}

        task = self.background_tasks[task_id]
        if not task.done():
            task.cancel()
            try:
                await task
            except asyncio.CancelledError:
                pass

        del self.background_tasks[task_id]
        return {"status": "cancelled"}

    async def cleanup(self) -> None:
        """Clean up all Neovim instances and background tasks."""
        for task_id, task in list(self.background_tasks.items()):
            if not task.done():
                task.cancel()
                try:
                    await task
                except asyncio.CancelledError:
                    pass
            del self.background_tasks[task_id]

        for instance_id in list(self.active_instances.keys()):
            await self.stop_instance(instance_id)

        self.logger.info(
            "Cleaned up all Neovim instances and background tasks"
        )


neovim_tool = NeovimTool()
