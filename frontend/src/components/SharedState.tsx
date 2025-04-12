import { useEffect, useState, useCallback } from 'react';

interface SharedStateProps {
  stateType: string;
  stateId: string;
  onStateChange?: (state: any) => void;
}

interface ServerAction {
  name: string;
  description?: string;
}

export const useSharedState = (stateType: string, stateId: string) => {
  const [state, setState] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchState = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(`/api/state/${stateType}/${stateId}`);
      
      if (!response.ok) {
        throw new Error(`Failed to fetch state: ${response.statusText}`);
      }
      
      const data = await response.json();
      setState(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, [stateType, stateId]);

  const updateState = useCallback(async (data: any) => {
    try {
      setLoading(true);
      const response = await fetch(`/api/state/${stateType}/${stateId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      
      if (!response.ok) {
        throw new Error(`Failed to update state: ${response.statusText}`);
      }
      
      await fetchState();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, [stateType, stateId, fetchState]);

  useEffect(() => {
    fetchState();
    
    const eventSource = new EventSource(`/api/events?state_type=${stateType}&state_id=${stateId}`);
    
    eventSource.addEventListener('update', (event) => {
      try {
        const data = JSON.parse(event.data);
        setState(data);
      } catch (err) {
        console.error('Error parsing SSE event:', err);
      }
    });
    
    eventSource.onerror = () => {
      console.error('SSE connection error');
      eventSource.close();
    };
    
    return () => {
      eventSource.close();
    };
  }, [stateType, stateId, fetchState]);

  return { state, loading, error, updateState };
};

export const useServerActions = () => {
  const [actions, setActions] = useState<ServerAction[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchActions = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/actions');
      
      if (!response.ok) {
        throw new Error(`Failed to fetch actions: ${response.statusText}`);
      }
      
      const data = await response.json();
      setActions(data.actions || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  const executeAction = useCallback(async (actionName: string, params: any = {}) => {
    try {
      const response = await fetch('/api/actions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          action: actionName,
          params,
        }),
      });
      
      if (!response.ok) {
        throw new Error(`Failed to execute action: ${response.statusText}`);
      }
      
      return await response.json();
    } catch (err) {
      throw err;
    }
  }, []);

  useEffect(() => {
    fetchActions();
  }, [fetchActions]);

  return { actions, loading, error, executeAction };
};

export const SharedStateComponent: React.FC<SharedStateProps> = ({ 
  stateType, 
  stateId,
  onStateChange 
}) => {
  const { state, loading, error, updateState } = useSharedState(stateType, stateId);
  
  useEffect(() => {
    if (state && onStateChange) {
      onStateChange(state);
    }
  }, [state, onStateChange]);
  
  if (loading) {
    return <div>Loading shared state...</div>;
  }
  
  if (error) {
    return <div>Error: {error}</div>;
  }
  
  if (!state) {
    return <div>No state found</div>;
  }
  
  return (
    <div className="shared-state">
      <h3>Shared State: {stateType}/{stateId}</h3>
      <pre>{JSON.stringify(state, null, 2)}</pre>
    </div>
  );
};

export default SharedStateComponent;
