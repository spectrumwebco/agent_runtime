"""
Shared State Module for Teemo
Handles state management and synchronization between frontend and backend
"""

from typing import Dict, List, Any, Optional
from pydantic import BaseModel, Field

class StateEvent(BaseModel):
    """Pydantic model for state events"""
    event_type: str = Field(..., description="Type of state event")
    source: str = Field(..., description="Source of the event")
    data: Dict[str, Any] = Field(default_factory=dict, description="Event data payload")
    timestamp: int = Field(..., description="Event timestamp")

class SharedStateModule:
    """
    Module for handling shared state between frontend and backend
    Supports state synchronization across React, Vue, and other frameworks
    """
    
    @property
    def name(self):
        """Returns the module name"""
        return "shared_state"
    
    @property
    def description(self):
        """Returns the module description"""
        return "Handles state management and synchronization between frontend and backend"
    
    @property
    def tools(self):
        """Returns a list of tools provided by the module"""
        return [
            "initialize_state",
            "update_state",
            "subscribe_to_state",
            "get_state",
            "sync_state"
        ]
    
    def initialize(self, context):
        """Initializes the module with execution context"""
        self.context = context
        self.framework = context.get("framework", "react")
        self.state_registry = {}
        self.subscribers = {}
        self.websocket_enabled = False
        self.sse_enabled = False
        return True
    
    def cleanup(self):
        """Cleans up module resources"""
        self.state_registry = {}
        self.subscribers = {}
        return True
    
    def initialize_state(self, state_spec):
        """
        Initializes shared state based on specifications
        
        Args:
            state_spec (dict): State specifications including:
                - name: State store name
                - initial_state: Initial state values
                - persistence: Whether state should persist
                - sync_method: Method for state synchronization (websocket, sse, etc.)
                - framework: Target framework (react, vue, etc.)
                
        Returns:
            dict: Initialized state configuration
        """
        framework = state_spec.get("framework", self.framework)
        sync_method = state_spec.get("sync_method", "websocket")
        
        if sync_method == "websocket":
            self.websocket_enabled = True
        elif sync_method == "sse":
            self.sse_enabled = True
            
        self.state_registry[state_spec["name"]] = {
            "current": state_spec.get("initial_state", {}),
            "persistence": state_spec.get("persistence", False),
            "sync_method": sync_method,
            "framework": framework
        }
        
        if framework == "react" or framework == "react-ts":
            code = self._generate_react_state_code(state_spec)
        elif framework == "vue":
            code = self._generate_vue_state_code(state_spec)
        else:
            code = f"// Generated state management code for {framework}"
            
        return {
            "code": code,
            "state_name": state_spec["name"],
            "sync_method": sync_method,
            "framework": framework
        }
    
    def update_state(self, state_name, updates, source="frontend"):
        """
        Updates shared state and notifies subscribers
        
        Args:
            state_name (str): Name of the state to update
            updates (dict): Updates to apply to the state
            source (str): Source of the update (frontend, backend)
            
        Returns:
            dict: Updated state
        """
        if state_name not in self.state_registry:
            return {"error": f"State '{state_name}' not found"}
            
        state = self.state_registry[state_name]
        state["current"].update(updates)
        
        event = StateEvent(
            event_type="state_update",
            source=source,
            data={"state_name": state_name, "updates": updates},
            timestamp=0  # Would use actual timestamp in implementation
        )
        
        if state_name in self.subscribers:
            for subscriber in self.subscribers[state_name]:
                subscriber(event)
                
        return {
            "state_name": state_name,
            "current_state": state["current"],
            "updated_fields": list(updates.keys())
        }
    
    def subscribe_to_state(self, state_name, callback):
        """
        Subscribes to state changes
        
        Args:
            state_name (str): Name of the state to subscribe to
            callback (callable): Function to call when state changes
            
        Returns:
            dict: Subscription information
        """
        if state_name not in self.state_registry:
            return {"error": f"State '{state_name}' not found"}
            
        if state_name not in self.subscribers:
            self.subscribers[state_name] = []
            
        self.subscribers[state_name].append(callback)
        
        return {
            "state_name": state_name,
            "subscription_id": len(self.subscribers[state_name]) - 1,
            "current_state": self.state_registry[state_name]["current"]
        }
    
    def get_state(self, state_name):
        """
        Gets current state
        
        Args:
            state_name (str): Name of the state to get
            
        Returns:
            dict: Current state
        """
        if state_name not in self.state_registry:
            return {"error": f"State '{state_name}' not found"}
            
        return {
            "state_name": state_name,
            "current_state": self.state_registry[state_name]["current"]
        }
    
    def sync_state(self, state_name, target="frontend"):
        """
        Synchronizes state between frontend and backend
        
        Args:
            state_name (str): Name of the state to synchronize
            target (str): Target to synchronize with (frontend, backend)
            
        Returns:
            dict: Synchronization result
        """
        if state_name not in self.state_registry:
            return {"error": f"State '{state_name}' not found"}
            
        state = self.state_registry[state_name]
        
        
        return {
            "state_name": state_name,
            "sync_method": state["sync_method"],
            "target": target,
            "current_state": state["current"]
        }
    
    def _generate_react_state_code(self, state_spec):
        """Generates React state management code using Jotai and Zustand"""
        name = state_spec["name"]
        initial_state = state_spec.get("initial_state", {})
        sync_method = state_spec.get("sync_method", "websocket")
        
        if "jotai" in state_spec.get("libraries", []):
            code = f"""// Generated Jotai state management for {name}
import {{ atom, useAtom }} from 'jotai';

// Define atoms
{self._generate_jotai_atoms(name, initial_state)}

// Define hooks
export const use{name.capitalize()}State = () => {{
  {self._generate_jotai_hooks(name, initial_state)}
  
  // Sync with backend
  {self._generate_sync_code(name, sync_method, "react")}
  
  return {{
    {", ".join([f"{key}" for key in initial_state.keys()])},
    {", ".join([f"set{key.capitalize()}" for key in initial_state.keys()])}
  }};
}};
"""
        else:
            code = f"""// Generated Zustand state management for {name}
import {{ create }} from 'zustand';

// Define store
export const use{name.capitalize()}Store = create((set) => ({{
  {self._generate_zustand_state(initial_state)},
  
  // Actions
  {self._generate_zustand_actions(initial_state)},
  
  // Sync with backend
  syncWithBackend: () => {{
    {self._generate_sync_code(name, sync_method, "react")}
  }}
}}));
"""
        
        return code
    
    def _generate_vue_state_code(self, state_spec):
        """Generates Vue state management code"""
        name = state_spec["name"]
        initial_state = state_spec.get("initial_state", {})
        sync_method = state_spec.get("sync_method", "websocket")
        
        code = f"""// Generated Vue state management for {name}
import {{ ref, reactive }} from 'vue';

export const use{name.capitalize()}State = () => {{
  // State
  const state = reactive({{
    {self._generate_vue_state(initial_state)}
  }});
  
  // Actions
  {self._generate_vue_actions(initial_state)}
  
  // Sync with backend
  {self._generate_sync_code(name, sync_method, "vue")}
  
  return {{
    ...state,
    {", ".join([f"update{key.capitalize()}" for key in initial_state.keys()])},
    syncWithBackend
  }};
}};
"""
        
        return code
    
    def _generate_jotai_atoms(self, name, initial_state):
        """Generates Jotai atom definitions"""
        atoms = []
        for key, value in initial_state.items():
            atoms.append(f"export const {key}Atom = atom({self._format_value(value)});")
        return "\n".join(atoms)
    
    def _generate_jotai_hooks(self, name, initial_state):
        """Generates Jotai hook usage"""
        hooks = []
        for key in initial_state.keys():
            hooks.append(f"const [{key}, set{key.capitalize()}] = useAtom({key}Atom);")
        return "\n  ".join(hooks)
    
    def _generate_zustand_state(self, initial_state):
        """Generates Zustand state definitions"""
        state = []
        for key, value in initial_state.items():
            state.append(f"{key}: {self._format_value(value)}")
        return ",\n  ".join(state)
    
    def _generate_zustand_actions(self, initial_state):
        """Generates Zustand action definitions"""
        actions = []
        for key in initial_state.keys():
            actions.append(f"set{key.capitalize()}: (value) => set({{ {key}: value }})")
        return ",\n  ".join(actions)
    
    def _generate_vue_state(self, initial_state):
        """Generates Vue state definitions"""
        state = []
        for key, value in initial_state.items():
            state.append(f"{key}: {self._format_value(value)}")
        return ",\n    ".join(state)
    
    def _generate_vue_actions(self, initial_state):
        """Generates Vue action definitions"""
        actions = []
        for key in initial_state.keys():
            actions.append(f"""const update{key.capitalize()} = (value) => {{
    state.{key} = value;
  }};""")
        return "\n  ".join(actions)
    
    def _generate_sync_code(self, name, sync_method, framework):
        """Generates code for state synchronization"""
        if sync_method == "websocket":
            return f"""// WebSocket sync implementation
const syncWithBackend = () => {{
  const socket = new WebSocket(`${{window.location.protocol === 'https:' ? 'wss:' : 'ws:'}}//${{window.location.host}}/api/state/{name}`);
  
  socket.onmessage = (event) => {{
    const data = JSON.parse(event.data);
    if (data.type === 'state_update') {{
      // Update local state with backend changes
      {self._generate_state_update_code(framework)}
    }}
  }};
  
  // Send local changes to backend
  {self._generate_state_send_code(framework)}
  
  return () => socket.close();
}};"""
        elif sync_method == "sse":
            return f"""// Server-Sent Events sync implementation
const syncWithBackend = () => {{
  const eventSource = new EventSource(`/api/state/{name}/events`);
  
  eventSource.addEventListener('state_update', (event) => {{
    const data = JSON.parse(event.data);
    // Update local state with backend changes
    {self._generate_state_update_code(framework)}
  }});
  
  // Send local changes to backend using fetch
  {self._generate_state_send_code(framework, "fetch")}
  
  return () => eventSource.close();
}};"""
        else:
            return f"""// Polling sync implementation
const syncWithBackend = () => {{
  const intervalId = setInterval(async () => {{
    try {{
      const response = await fetch(`/api/state/{name}`);
      const data = await response.json();
      // Update local state with backend changes
      {self._generate_state_update_code(framework)}
    }} catch (error) {{
      console.error('Failed to sync state:', error);
    }}
  }}, 5000); // Poll every 5 seconds
  
  return () => clearInterval(intervalId);
}};"""
    
    def _generate_state_update_code(self, framework):
        """Generates code for updating state from backend"""
        if framework == "react":
            return """Object.entries(data.updates).forEach(([key, value]) => {
        // For Zustand
        set({ [key]: value });
        // For Jotai (would need to be customized per atom)
      });"""
        elif framework == "vue":
            return """Object.entries(data.updates).forEach(([key, value]) => {
        state[key] = value;
      });"""
        else:
            return "// Update state based on framework"
    
    def _generate_state_send_code(self, framework, method="websocket"):
        """Generates code for sending state to backend"""
        if method == "websocket":
            return """// Send updates to backend
    const sendUpdate = (updates) => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
          type: 'state_update',
          updates
        }));
      }
    };"""
        else:
            return """// Send updates to backend
    const sendUpdate = async (updates) => {
      try {
        await fetch(`/api/state/${name}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            type: 'state_update',
            updates
          })
        });
      } catch (error) {
        console.error('Failed to send state update:', error);
      }
    };"""
    
    def _format_value(self, value):
        """Formats a value for code generation"""
        if isinstance(value, str):
            return f"'{value}'"
        elif isinstance(value, bool):
            return str(value).lower()
        elif isinstance(value, (int, float)):
            return str(value)
        elif isinstance(value, dict):
            items = [f"{k}: {self._format_value(v)}" for k, v in value.items()]
            return "{ " + ", ".join(items) + " }"
        elif isinstance(value, list):
            items = [self._format_value(item) for item in value]
            return "[" + ", ".join(items) + "]"
        elif value is None:
            return "null"
        else:
            return str(value)
