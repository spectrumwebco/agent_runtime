import React from 'react';
import { ThemeProvider } from '../components/theme-provider';
import { QueryClientProvider } from './query-client-provider';

interface RootProviderProps {
  children: React.ReactNode;
}

export function RootProvider({ children }: RootProviderProps) {
  return (
    <ThemeProvider>
      <QueryClientProvider>
        {children}
      </QueryClientProvider>
    </ThemeProvider>
  );
}
