import React, { ReactNode } from 'react';
import { PlatformOnly } from './platform-specific';

interface ElectronWrapperProps {
  children: ReactNode;
  menuTemplate?: any[]; // Electron menu template
  windowOptions?: {
    width?: number;
    height?: number;
    minWidth?: number;
    minHeight?: number;
    title?: string;
    icon?: string;
  };
}

/**
 * Wrapper component for Electron-specific functionality
 * Only renders on Electron platform
 */
export const ElectronWrapper: React.FC<ElectronWrapperProps> = ({
  children,
  menuTemplate,
  windowOptions,
}) => {
  React.useEffect(() => {
    if (typeof window !== 'undefined' && 
        window.electron) {
      
      if (menuTemplate) {
        window.electron.setApplicationMenu(menuTemplate);
      }
      
      if (windowOptions) {
        window.electron.configureWindow(windowOptions);
      }
      
      const handleBeforeUnload = (e: BeforeUnloadEvent) => {
        e.preventDefault();
        e.returnValue = '';
        
        window.electron.confirmClose();
      };
      
      window.addEventListener('beforeunload', handleBeforeUnload);
      
      return () => {
        window.removeEventListener('beforeunload', handleBeforeUnload);
      };
    }
  }, [menuTemplate, windowOptions]);
  
  return (
    <PlatformOnly platform="electron">
      {children}
    </PlatformOnly>
  );
};

export default ElectronWrapper;
