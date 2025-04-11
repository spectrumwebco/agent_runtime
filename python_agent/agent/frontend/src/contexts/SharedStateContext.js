import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import SharedStateClient from '../utils/sharedState';

const DEFAULT_SERVER_URL = 'ws://localhost:8080/ws';

const SharedStateContext = createContext(null);

/**
 * Provider component for shared state
 * @param {object} props - Component props
 * @param {string} props.serverUrl - WebSocket server URL
 * @param {React.ReactNode} props.children - Child components
 * @returns {React.ReactElement} - Provider component
 */
export const SharedStateProvider = ({ serverUrl = DEFAULT_SERVER_URL, children }) => {
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState(null);
  const clientRef = useRef(null);
  const stateCache = useRef(new Map());

  useEffect(() => {
    if (!clientRef.current) {
      clientRef.current = new SharedStateClient(serverUrl);
    }

    const client = clientRef.current;

    client.connect()
      .then(() => {
        setIsConnected(true);
        setError(null);
      })
      .catch(err => {
        console.error('Error connecting to shared state server:', err);
        setError(err);
      });

    return () => {
      if (client) {
        client.close();
      }
    };
  }, [serverUrl]);

  /**
   * Subscribe to state updates
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {function} callback - The callback to call when state is updated
   */
  const subscribeToState = (stateType, stateId, callback) => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    clientRef.current.subscribeToState(stateType, stateId, callback);
  };

  /**
   * Unsubscribe from state updates
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {function} callback - The callback to remove
   */
  const unsubscribeFromState = (stateType, stateId, callback) => {
    if (!clientRef.current) {
      return;
    }

    clientRef.current.unsubscribeFromState(stateType, stateId, callback);
  };

  /**
   * Update state
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {object} data - The state data
   */
  const updateState = async (stateType, stateId, data) => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    if (!isConnected) {
      await clientRef.current.connect();
      setIsConnected(true);
    }

    const cacheKey = `${stateType}:${stateId}`;
    const currentState = stateCache.current.get(cacheKey) || {};
    const newState = typeof data === 'function' ? data(currentState) : data;
    
    const mergedState = typeof newState === 'object' && !Array.isArray(newState)
      ? { ...currentState, ...newState }
      : newState;
    
    stateCache.current.set(cacheKey, mergedState);

    clientRef.current.updateState(stateType, stateId, mergedState);
  };

  /**
   * Send an event to the server
   * @param {object} eventData - The event data
   */
  const sendEvent = async (eventData) => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    if (!isConnected) {
      await clientRef.current.connect();
      setIsConnected(true);
    }

    clientRef.current.sendEvent(eventData);
  };

  const value = {
    client: clientRef.current,
    isConnected,
    error,
    subscribeToState,
    unsubscribeFromState,
    updateState,
    sendEvent,
  };

  return (
    <SharedStateContext.Provider value={value}>
      {children}
    </SharedStateContext.Provider>
  );
};

/**
 * Hook for using the shared state context
 * @returns {object} - Shared state context
 */
export const useSharedStateContext = () => {
  const context = useContext(SharedStateContext);
  if (!context) {
    throw new Error('useSharedStateContext must be used within a SharedStateProvider');
  }
  return context;
};

export default SharedStateContext;
