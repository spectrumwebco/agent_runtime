import React, { createContext, useContext, useEffect, useRef, useState, ReactNode } from 'react';
import SharedStateClient from '../utils/sharedState';

const DEFAULT_SERVER_URL = 'ws://localhost:8080/ws';

interface SharedStateContextValue {
  client: SharedStateClient | null;
  isConnected: boolean;
  error: Error | null;
  subscribeToState: (stateType: string, stateId: string, callback: StateCallback) => void;
  unsubscribeFromState: (stateType: string, stateId: string, callback?: StateCallback) => void;
  updateState: (stateType: string, stateId: string, data: StateUpdateParam) => Promise<void>;
  sendEvent: (eventData: Record<string, any>) => Promise<void>;
}

type StateCallback = (data: any) => void;
type StateUpdateFunction = (prevState: Record<string, any>) => Record<string, any>;
type StateUpdateParam = Record<string, any> | StateUpdateFunction;

const SharedStateContext = createContext<SharedStateContextValue | null>(null);

interface SharedStateProviderProps {
  serverUrl?: string;
  children: ReactNode;
}

/**
 * Provider component for shared state
 * @param props - Component props
 * @param props.serverUrl - WebSocket server URL
 * @param props.children - Child components
 * @returns Provider component
 */
export const SharedStateProvider: React.FC<SharedStateProviderProps> = ({ 
  serverUrl = DEFAULT_SERVER_URL, 
  children 
}) => {
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const clientRef = useRef<SharedStateClient | null>(null);
  const stateCache = useRef<Map<string, Record<string, any>>>(new Map());

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
        setError(err instanceof Error ? err : new Error(String(err)));
      });

    return () => {
      if (client) {
        client.close();
      }
    };
  }, [serverUrl]);

  /**
   * Subscribe to state updates
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param callback - The callback to call when state is updated
   */
  const subscribeToState = (stateType: string, stateId: string, callback: StateCallback): void => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    clientRef.current.subscribeToState(stateType, stateId, callback);
  };

  /**
   * Unsubscribe from state updates
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param callback - The callback to remove
   */
  const unsubscribeFromState = (stateType: string, stateId: string, callback?: StateCallback): void => {
    if (!clientRef.current) {
      return;
    }

    clientRef.current.unsubscribeFromState(stateType, stateId, callback || null);
  };

  /**
   * Update state
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param data - The state data
   */
  const updateState = async (stateType: string, stateId: string, data: StateUpdateParam): Promise<void> => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    if (!isConnected) {
      await clientRef.current.connect();
      setIsConnected(true);
    }

    const cacheKey = `${stateType}:${stateId}`;
    const currentState = stateCache.current.get(cacheKey) || {};
    const newState = typeof data === 'function' 
      ? (data as StateUpdateFunction)(currentState) 
      : data;
    
    const mergedState = typeof newState === 'object' && !Array.isArray(newState)
      ? { ...currentState, ...newState }
      : newState;
    
    stateCache.current.set(cacheKey, mergedState);

    clientRef.current.updateState(stateType, stateId, mergedState);
  };

  /**
   * Send an event to the server
   * @param eventData - The event data
   */
  const sendEvent = async (eventData: Record<string, any>): Promise<void> => {
    if (!clientRef.current) {
      throw new Error('Client not initialized');
    }

    if (!isConnected) {
      await clientRef.current.connect();
      setIsConnected(true);
    }

    clientRef.current.sendEvent(eventData);
  };

  const value: SharedStateContextValue = {
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
 * @returns Shared state context
 */
export const useSharedStateContext = (): SharedStateContextValue => {
  const context = useContext(SharedStateContext);
  if (!context) {
    throw new Error('useSharedStateContext must be used within a SharedStateProvider');
  }
  return context;
};

export default SharedStateContext;
