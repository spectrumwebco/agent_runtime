import React, { createContext, useContext, ReactNode } from 'react';
import { Platform, currentPlatform, platformFeatures } from '../../utils/platform';

interface PlatformContextType {
  platform: Platform;
  isWeb: boolean;
  isElectron: boolean;
  isMobile: boolean;
  features: typeof platformFeatures;
}

const PlatformContext = createContext<PlatformContextType>({
  platform: 'web',
  isWeb: true,
  isElectron: false,
  isMobile: false,
  features: platformFeatures,
});

export const usePlatform = () => useContext(PlatformContext);

interface PlatformProviderProps {
  children: ReactNode;
  forcePlatform?: Platform;
}

export const PlatformProvider: React.FC<PlatformProviderProps> = ({
  children,
  forcePlatform,
}) => {
  const platform = forcePlatform || currentPlatform;
  
  const value = {
    platform,
    isWeb: platform === 'web',
    isElectron: platform === 'electron',
    isMobile: platform === 'mobile',
    features: platformFeatures,
  };
  
  return (
    <PlatformContext.Provider value={value}>
      {children}
    </PlatformContext.Provider>
  );
};

export default PlatformProvider;
