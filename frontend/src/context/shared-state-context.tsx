import React, { createContext, useContext, useState, useEffect } from 'react';

export interface SharedState {
  agentState: {
    status: 'idle' | 'running' | 'completed' | 'error';
    progress: number;
    currentTask: string;
    completedTasks: string[];
    pendingTasks: string[];
  };
  uiState: {
    theme: 'light' | 'dark';
    layout: 'default' | 'flipped';
    sidebarOpen: boolean;
  };
  modelState: {
    selectedModel: string;
    availableModels: string[];
  };
}

const defaultState: SharedState = {
  agentState: {
    status: 'idle',
    progress: 0,
    currentTask: '',
    completedTasks: [],
    pendingTasks: [],
  },
  uiState: {
    theme: 'dark',
    layout: 'flipped',
    sidebarOpen: true,
  },
  modelState: {
    selectedModel: 'gemini-2.5-pro',
    availableModels: ['gemini-2.5-pro', 'llama-4', 'gpt-4o'],
  },
};

interface SharedStateContextType {
  state: SharedState;
  updateAgentState: (agentState: Partial<SharedState['agentState']>) => void;
  updateUIState: (uiState: Partial<SharedState['uiState']>) => void;
  updateModelState: (modelState: Partial<SharedState['modelState']>) => void;
}

const SharedStateContext = createContext<SharedStateContextType | undefined>(undefined);

export const SharedStateProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, setState] = useState<SharedState>(defaultState);

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8000/ws/state');

    socket.onopen = () => {
      console.log('Connected to shared state server');
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setState((prevState) => ({
          ...prevState,
          ...data,
        }));
      } catch (error) {
        console.error('Error parsing shared state update:', error);
      }
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socket.onclose = () => {
      console.log('Disconnected from shared state server');
    };

    return () => {
      socket.close();
    };
  }, []);

  const updateAgentState = (agentState: Partial<SharedState['agentState']>) => {
    setState((prevState) => ({
      ...prevState,
      agentState: {
        ...prevState.agentState,
        ...agentState,
      },
    }));
  };

  const updateUIState = (uiState: Partial<SharedState['uiState']>) => {
    setState((prevState) => ({
      ...prevState,
      uiState: {
        ...prevState.uiState,
        ...uiState,
      },
    }));
  };

  const updateModelState = (modelState: Partial<SharedState['modelState']>) => {
    setState((prevState) => ({
      ...prevState,
      modelState: {
        ...prevState.modelState,
        ...modelState,
      },
    }));
  };

  return (
    <SharedStateContext.Provider
      value={{
        state,
        updateAgentState,
        updateUIState,
        updateModelState,
      }}
    >
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
