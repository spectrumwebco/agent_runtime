import React, { createContext, useContext, useState, useEffect } from "react";

export interface SharedAppState {
  activeAgentId: string | null;
  activeModelId: string | null;
  activeToolIds: string[];
  activeView: "code" | "terminal" | "browser" | "graph" | "chat";
  viewMode: "control" | "worker" | null;
  isGenerating: boolean;
  progress: number;
  lastAction: string | null;
  lastError: string | null;
}

const defaultState: SharedAppState = {
  activeAgentId: null,
  activeModelId: null,
  activeToolIds: [],
  activeView: "code",
  viewMode: "control",
  isGenerating: false,
  progress: 0,
  lastAction: null,
  lastError: null,
};

interface SharedStateContextType {
  state: SharedAppState;
  updateState: (updates: Partial<SharedAppState>) => void;
  resetState: () => void;
}

const SharedStateContext = createContext<SharedStateContextType | undefined>(undefined);

interface SharedStateProviderProps {
  children: React.ReactNode;
  initialState?: Partial<SharedAppState>;
}

export const SharedStateProvider: React.FC<SharedStateProviderProps> = ({
  children,
  initialState = {},
}) => {
  const [state, setState] = useState<SharedAppState>({
    ...defaultState,
    ...initialState,
  });

  useEffect(() => {
    const ws = new WebSocket(`${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws/state`);
    
    ws.onopen = () => {
      console.log('WebSocket connection established');
    };
    
    ws.onmessage = (event) => {
      try {
        const updates = JSON.parse(event.data);
        setState(prevState => ({
          ...prevState,
          ...updates,
        }));
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };
    
    return () => {
      ws.close();
    };
  }, []);

  const updateState = (updates: Partial<SharedAppState>) => {
    setState(prevState => ({
      ...prevState,
      ...updates,
    }));
    
    fetch('/api/state', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updates),
    }).catch(error => {
      console.error('Error sending state update to backend:', error);
    });
  };

  const resetState = () => {
    setState(defaultState);
    
    fetch('/api/state/reset', {
      method: 'POST',
    }).catch(error => {
      console.error('Error resetting state on backend:', error);
    });
  };

  return (
    <SharedStateContext.Provider value={{ state, updateState, resetState }}>
      {children}
    </SharedStateContext.Provider>
  );
};

export const useSharedState = () => {
  const context = useContext(SharedStateContext);
  
  if (context === undefined) {
    throw new Error('useSharedState must be used within a SharedStateProvider');
  }
  
  return context;
};
