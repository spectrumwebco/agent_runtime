import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

export type Theme = 'light' | 'dark' | 'system';

export enum AgentState {
  IDLE = 'idle',
  THINKING = 'thinking',
  EXECUTING = 'executing',
  ERROR = 'error',
  COMPLETED = 'completed',
}

interface UIState {
  theme: Theme;
  sidebarOpen: boolean;
  setTheme: (theme: Theme) => void;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
}

interface AgentStateStore {
  currentState: AgentState;
  lastError: string | null;
  setCurrentState: (state: AgentState) => void;
  setLastError: (error: string | null) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      theme: 'system',
      sidebarOpen: true,
      setTheme: (theme) => set({ theme }),
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
    }),
    {
      name: 'kled-ui-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export const useAgentStore = create<AgentStateStore>((set) => ({
  currentState: AgentState.IDLE,
  lastError: null,
  setCurrentState: (state) => set({ currentState: state }),
  setLastError: (error) => set({ lastError: error }),
}));

interface UserSettings {
  aiModel: string;
  githubToken: string | null;
  settings: Record<string, any>;
  setAIModel: (model: string) => void;
  setGithubToken: (token: string | null) => void;
  updateSettings: (key: string, value: any) => void;
}

export const useUserSettings = create<UserSettings>()(
  persist(
    (set) => ({
      aiModel: 'llama-4',
      githubToken: null,
      settings: {},
      setAIModel: (model) => set({ aiModel: model }),
      setGithubToken: (token) => set({ githubToken: token }),
      updateSettings: (key, value) =>
        set((state) => ({
          settings: {
            ...state.settings,
            [key]: value,
          },
        })),
    }),
    {
      name: 'kled-user-settings',
      storage: createJSONStorage(() => localStorage),
    }
  )
);
