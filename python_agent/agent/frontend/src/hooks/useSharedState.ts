import { useState, useEffect, useRef, useCallback } from 'react';
import SharedStateClient from '../utils/sharedState';

const DEFAULT_SERVER_URL = 'ws://localhost:8080/ws';

type StateData = Record<string, any> | null;
type StateUpdateFunction = (prevState: StateData) => StateData;
type StateUpdateParam = StateData | StateUpdateFunction;

interface SharedStateHookResult {
  state: StateData;
  loading: boolean;
  error: Error | null;
  updateState: (newState: StateUpdateParam) => Promise<void>;
  sendEvent: (eventData: Record<string, any>) => Promise<void>;
  client: SharedStateClient | null;
}

/**
 * Hook for using shared state
 * @param stateType - The type of state
 * @param stateId - The ID of the state
 * @param initialState - Initial state value
 * @param serverUrl - WebSocket server URL
 * @returns State, loading status, error, and update function
 */
export const useSharedState = (
  stateType: string,
  stateId: string,
  initialState: StateData = null,
  serverUrl: string = DEFAULT_SERVER_URL
): SharedStateHookResult => {
  const [state, setState] = useState<StateData>(initialState);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  const clientRef = useRef<SharedStateClient | null>(null);
  const stateTypeRef = useRef<string>(stateType);
  const stateIdRef = useRef<string>(stateId);

  useEffect(() => {
    stateTypeRef.current = stateType;
    stateIdRef.current = stateId;
  }, [stateType, stateId]);

  useEffect(() => {
    if (!clientRef.current) {
      clientRef.current = new SharedStateClient(serverUrl);
    }

    const client = clientRef.current;
    let mounted = true;

    const handleStateUpdate = (newState: any) => {
      if (mounted) {
        setState(prevState => {
          if (typeof newState === 'function') {
            return newState(prevState);
          }
          if (newState && typeof newState === 'object' && !Array.isArray(newState)) {
            return { ...prevState, ...newState };
          }
          return newState;
        });
        setLoading(false);
      }
    };

    const fetchInitialState = async () => {
      try {
        await client.connect();

        client.subscribeToState(stateType, stateId, handleStateUpdate);

        setTimeout(() => {
          if (mounted && loading) {
            setLoading(false);
          }
        }, 2000);
      } catch (err) {
        if (mounted) {
          console.error('Error connecting to shared state:', err);
          setError(err instanceof Error ? err : new Error(String(err)));
          setLoading(false);
        }
      }
    };

    fetchInitialState();

    return () => {
      mounted = false;
      if (client) {
        client.unsubscribeFromState(stateType, stateId);
      }
    };
  }, [stateType, stateId, serverUrl, loading]);

  const updateState = useCallback(async (newState: StateUpdateParam) => {
    try {
      if (!clientRef.current) {
        throw new Error('Client not initialized');
      }

      await clientRef.current.connect();

      clientRef.current.updateState(
        stateTypeRef.current,
        stateIdRef.current,
        typeof newState === 'function' 
          ? newState(state) as Record<string, any>
          : newState as Record<string, any>
      );

      setState(prevState => {
        if (typeof newState === 'function') {
          return newState(prevState);
        }
        if (newState && typeof newState === 'object' && !Array.isArray(newState)) {
          return { ...prevState, ...newState };
        }
        return newState;
      });
    } catch (err) {
      console.error('Error updating state:', err);
      setError(err instanceof Error ? err : new Error(String(err)));
      throw err;
    }
  }, [state]);

  const sendEvent = useCallback(async (eventData: Record<string, any>) => {
    try {
      if (!clientRef.current) {
        throw new Error('Client not initialized');
      }

      await clientRef.current.connect();

      clientRef.current.sendEvent(eventData);
    } catch (err) {
      console.error('Error sending event:', err);
      setError(err instanceof Error ? err : new Error(String(err)));
      throw err;
    }
  }, []);

  return {
    state,
    loading,
    error,
    updateState,
    sendEvent,
    client: clientRef.current,
  };
};

export default useSharedState;
