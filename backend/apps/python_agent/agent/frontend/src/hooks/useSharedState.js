import { useState, useEffect, useRef, useCallback } from 'react';
import SharedStateClient from '../utils/sharedState';

const DEFAULT_SERVER_URL = 'ws://localhost:8080/ws';

/**
 * Hook for using shared state
 * @param {string} stateType - The type of state
 * @param {string} stateId - The ID of the state
 * @param {any} initialState - Initial state value
 * @param {string} serverUrl - WebSocket server URL
 * @returns {object} - State, loading status, error, and update function
 */
export const useSharedState = (stateType, stateId, initialState = null, serverUrl = DEFAULT_SERVER_URL) => {
  const [state, setState] = useState(initialState);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const clientRef = useRef(null);
  const stateTypeRef = useRef(stateType);
  const stateIdRef = useRef(stateId);

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

    const handleStateUpdate = (newState) => {
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
          setError(err);
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

  const updateState = useCallback(async (newState) => {
    try {
      if (!clientRef.current) {
        throw new Error('Client not initialized');
      }

      await clientRef.current.connect();

      clientRef.current.updateState(
        stateTypeRef.current,
        stateIdRef.current,
        newState
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
      setError(err);
      throw err;
    }
  }, []);

  const sendEvent = useCallback(async (eventData) => {
    try {
      if (!clientRef.current) {
        throw new Error('Client not initialized');
      }

      await clientRef.current.connect();

      clientRef.current.sendEvent(eventData);
    } catch (err) {
      console.error('Error sending event:', err);
      setError(err);
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
