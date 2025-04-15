import { useState, useEffect, useCallback } from 'react';

interface UseSharedStateOptions {
  stateType?: string;
  stateId?: string;
  initialState?: any;
}

/**
 * Custom hook for using shared application state between
 * React frontend and Go/Python backend.
 * 
 * @param options Configuration options for the shared state
 * @returns Object containing state, update function, and connection status
 */
const useSharedState = ({
  stateType = 'shared',
  stateId = 'default',
  initialState = {}
}: UseSharedStateOptions = {}) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [state, setState] = useState(initialState);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = process.env.REACT_APP_API_HOST || window.location.host;
    const wsUrl = `${protocol}//${host}/ws/state/${stateType}/${stateId}/`;
    
    console.log(`Connecting to WebSocket at ${wsUrl}`);
    
    const newSocket = new WebSocket(wsUrl);
    
    newSocket.onopen = () => {
      console.log('WebSocket connection established');
      setConnected(true);
      setError(null);
      
      newSocket.send(JSON.stringify({
        type: 'get_state'
      }));
    };
    
    newSocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        if (data.type === 'state_update') {
          console.log('Received state update:', data);
          setState(data.data);
        }
      } catch (err) {
        console.error('Error parsing WebSocket message:', err);
      }
    };
    
    newSocket.onerror = (event) => {
      console.error('WebSocket error:', event);
      setError('WebSocket connection error');
    };
    
    newSocket.onclose = () => {
      console.log('WebSocket connection closed');
      setConnected(false);
    };
    
    setSocket(newSocket);
    
    return () => {
      if (newSocket) {
        newSocket.close();
      }
    };
  }, [stateType, stateId]);
  
  const updateState = useCallback((newState: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({
        type: 'update_state',
        data: newState
      }));
    } else {
      setError('WebSocket not connected');
    }
  }, [socket]);
  
  return {
    state,
    updateState,
    connected,
    error
  };
};

export default useSharedState;
