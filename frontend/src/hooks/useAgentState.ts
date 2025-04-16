import { useState, useEffect } from 'react';

interface AgentState {
  id: string;
  name: string;
  status: 'idle' | 'running' | 'paused' | 'completed' | 'failed';
  progress: number;
  current_task?: string;
  metadata?: Record<string, any>;
}

export function useAgentState(agentId: string) {
  const [agentState, setAgentState] = useState<AgentState | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchAgentState = async () => {
    try {
      const response = await fetch(`/api/agents/${agentId}/state`);
      
      if (!response.ok) {
        throw new Error(`Failed to fetch agent state: ${response.statusText}`);
      }
      
      const data = await response.json();
      setAgentState(data);
      setError(null);
    } catch (err) {
      console.error('Error fetching agent state:', err);
      setError(err instanceof Error ? err : new Error(String(err)));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAgentState();
    
    const intervalId = setInterval(fetchAgentState, 5000);
    
    return () => clearInterval(intervalId);
  }, [agentId]);

  const updateAgentState = async (updates: Partial<AgentState>) => {
    try {
      const response = await fetch(`/api/agents/${agentId}/state`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updates),
      });
      
      if (!response.ok) {
        throw new Error(`Failed to update agent state: ${response.statusText}`);
      }
      
      const data = await response.json();
      setAgentState(data);
      return data;
    } catch (err) {
      console.error('Error updating agent state:', err);
      setError(err instanceof Error ? err : new Error(String(err)));
      throw err;
    }
  };

  return { agentState, isLoading, error, updateAgentState, refreshAgentState: fetchAgentState };
}
