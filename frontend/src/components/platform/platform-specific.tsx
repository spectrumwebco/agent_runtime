import React, { ReactNode } from 'react';
import { usePlatform } from './platform-provider';

interface PlatformSpecificProps {
  web?: ReactNode;
  electron?: ReactNode;
  mobile?: ReactNode;
  fallback?: ReactNode;
}

/**
 * Component that renders different content based on the current platform
 * Use this for platform-specific UI elements
 */
export const PlatformSpecific: React.FC<PlatformSpecificProps> = ({
  web,
  electron,
  mobile,
  fallback,
}) => {
  const { isWeb, isElectron, isMobile } = usePlatform();
  
  if (isWeb && web !== undefined) {
    return <>{web}</>;
  }
  
  if (isElectron && electron !== undefined) {
    return <>{electron}</>;
  }
  
  if (isMobile && mobile !== undefined) {
    return <>{mobile}</>;
  }
  
  if (fallback !== undefined) {
    return <>{fallback}</>;
  }
  
  return null;
};

interface PlatformOnlyProps {
  platform: 'web' | 'electron' | 'mobile';
  children: ReactNode;
}

/**
 * Component that only renders its children on the specified platform
 */
export const PlatformOnly: React.FC<PlatformOnlyProps> = ({
  platform,
  children,
}) => {
  const { isWeb, isElectron, isMobile } = usePlatform();
  
  if (
    (platform === 'web' && isWeb) ||
    (platform === 'electron' && isElectron) ||
    (platform === 'mobile' && isMobile)
  ) {
    return <>{children}</>;
  }
  
  return null;
};

export default PlatformSpecific;
