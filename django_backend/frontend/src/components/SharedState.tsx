import React, { useEffect, useState, useCallback } from 'react';

interface SharedStateProps {
  stateType?: string;
  stateId?: string;
  onStateChange?: (data: any) => void;
  initialState?: any;
}

/**
 * SharedState component for real-time state synchronization between
 * React frontend and Go/Python backend.
 * 
 * This component establishes a WebSocket connection to the shared state
 * endpoint and provides real-time updates to the UI.
 */
const SharedState: React.FC<SharedStateProps> = ({
  stateType = 'shared',
  stateId = 'default',
  onStateChange,
  initialState = {}
}) => {
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
          
          if (onStateChange) {
            onStateChange(data.data);
          }
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
  }, [stateType, stateId, onStateChange]);
  
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

export default SharedState;
