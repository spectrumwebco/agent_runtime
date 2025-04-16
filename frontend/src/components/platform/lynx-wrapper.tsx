import React, { ReactNode, useEffect } from 'react';
import { PlatformOnly } from './platform-specific';

interface LynxWrapperProps {
  children: ReactNode;
  mobileOptions?: {
    statusBarStyle?: 'default' | 'light-content' | 'dark-content';
    navigationBarColor?: string;
    allowsBackForwardNavigationGestures?: boolean;
    screenOrientation?: 'portrait' | 'landscape' | 'auto';
  };
}

/**
 * Wrapper component for Lynx-React mobile-specific functionality
 * Only renders on mobile platform
 */
export const LynxWrapper: React.FC<LynxWrapperProps> = ({
  children,
  mobileOptions,
}) => {
  useEffect(() => {
    if (typeof window !== 'undefined' && 
        window.lynxReact) {
      
      if (mobileOptions) {
        if (window.lynxReact.setStatusBarStyle && mobileOptions.statusBarStyle) {
          window.lynxReact.setStatusBarStyle(mobileOptions.statusBarStyle);
        }
        
        if (window.lynxReact.setNavigationBarColor && mobileOptions.navigationBarColor) {
          window.lynxReact.setNavigationBarColor(mobileOptions.navigationBarColor);
        }
        
        if (window.lynxReact.setAllowsBackForwardNavigationGestures && 
            mobileOptions.allowsBackForwardNavigationGestures !== undefined) {
          window.lynxReact.setAllowsBackForwardNavigationGestures(
            mobileOptions.allowsBackForwardNavigationGestures
          );
        }
        
        if (window.lynxReact.setScreenOrientation && mobileOptions.screenOrientation) {
          window.lynxReact.setScreenOrientation(mobileOptions.screenOrientation);
        }
      }
      
      const handleAppStateChange = (state: string) => {
        console.log('Mobile app state changed:', state);
      };
      
      if (window.lynxReact.addEventListener) {
        window.lynxReact.addEventListener('appStateChange', handleAppStateChange);
        
        return () => {
          window.lynxReact.removeEventListener('appStateChange', handleAppStateChange);
        };
      }
    }
  }, [mobileOptions]);
  
  return (
    <PlatformOnly platform="mobile">
      {children}
    </PlatformOnly>
  );
};

export default LynxWrapper;
