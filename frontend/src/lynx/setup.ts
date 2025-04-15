import { LynxConfig, initializeLynx } from 'lynx-react';

/**
 * Configuration for Lynx React mobile integration
 */
export const lynxConfig: LynxConfig = {
  appName: 'Agent Runtime',
  appId: 'com.spectrumwebco.agentruntime',
  version: '1.0.0',
  theme: {
    primaryColor: '#10b981', // Emerald-500
    darkMode: {
      background: '#1a1a1a', // Charcoal grey
      text: '#ffffff',
    },
    lightMode: {
      background: '#f9fafb', // Lightest grey in Tailwind
      text: '#111827',
    },
  },
  permissions: {
    camera: false,
    location: false,
    notifications: true,
    storage: true,
  },
  deepLinking: {
    scheme: 'agentruntime',
    host: 'app',
  },
  platforms: {
    ios: {
      minVersion: '14.0',
      capabilities: [
        'BackgroundModes',
        'Push',
      ],
    },
    android: {
      minSdkVersion: 24,
      targetSdkVersion: 33,
      permissions: [
        'android.permission.INTERNET',
        'android.permission.ACCESS_NETWORK_STATE',
      ],
    },
  },
};

/**
 * Initialize Lynx React for mobile applications
 */
export const initializeMobileApp = () => {
  return initializeLynx(lynxConfig);
};

/**
 * Check if the app is running in a mobile environment
 */
export const isMobileApp = (): boolean => {
  return typeof window !== 'undefined' && 
    (window.navigator.userAgent.includes('LynxReact') || 
     /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(window.navigator.userAgent));
};

/**
 * Get the current platform (ios, android, or web)
 */
export const getPlatform = (): 'ios' | 'android' | 'web' => {
  if (typeof window === 'undefined') return 'web';
  
  const userAgent = window.navigator.userAgent.toLowerCase();
  
  if (userAgent.includes('iphone') || userAgent.includes('ipad') || userAgent.includes('ipod')) {
    return 'ios';
  }
  
  if (userAgent.includes('android')) {
    return 'android';
  }
  
  return 'web';
};
