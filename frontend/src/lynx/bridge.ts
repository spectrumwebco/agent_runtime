/**
 * Bridge between web app and native mobile functionality
 * Provides a consistent API for accessing native features
 */

export interface NativeBridge {
  readFile: (path: string) => Promise<string>;
  writeFile: (path: string, content: string) => Promise<boolean>;
  listFiles: (directory: string) => Promise<string[]>;
  
  getDeviceInfo: () => Promise<{
    platform: string;
    osVersion: string;
    appVersion: string;
    deviceId: string;
    deviceName: string;
  }>;
  
  scheduleNotification: (title: string, body: string, delay: number) => Promise<string>;
  cancelNotification: (id: string) => Promise<boolean>;
  
  shareText: (text: string, title?: string) => Promise<boolean>;
  shareFile: (path: string, title?: string) => Promise<boolean>;
  
  biometricAuth: (reason: string) => Promise<boolean>;
  
  getNetworkStatus: () => Promise<{
    connected: boolean;
    type: 'wifi' | 'cellular' | 'unknown' | 'none';
  }>;
}

const webBridge: NativeBridge = {
  readFile: async () => {
    throw new Error('File system access not available in web environment');
  },
  writeFile: async () => {
    throw new Error('File system access not available in web environment');
  },
  listFiles: async () => {
    throw new Error('File system access not available in web environment');
  },
  getDeviceInfo: async () => ({
    platform: 'web',
    osVersion: 'unknown',
    appVersion: '1.0.0',
    deviceId: 'web-browser',
    deviceName: 'Web Browser',
  }),
  scheduleNotification: async () => {
    console.warn('Notifications not available in web environment');
    return 'mock-id';
  },
  cancelNotification: async () => {
    console.warn('Notifications not available in web environment');
    return false;
  },
  shareText: async (text) => {
    if (navigator.share) {
      try {
        await navigator.share({ text });
        return true;
      } catch (error) {
        console.error('Error sharing text:', error);
        return false;
      }
    }
    console.warn('Web Share API not available');
    return false;
  },
  shareFile: async () => {
    console.warn('File sharing not available in web environment');
    return false;
  },
  biometricAuth: async () => {
    console.warn('Biometric authentication not available in web environment');
    return false;
  },
  getNetworkStatus: async () => {
    const online = typeof navigator !== 'undefined' ? navigator.onLine : false;
    return {
      connected: online,
      type: online ? 'unknown' : 'none',
    };
  },
};

export const getNativeBridge = (): NativeBridge => {
  if (typeof window !== 'undefined' && window.lynxBridge) {
    return window.lynxBridge as NativeBridge;
  }
  
  return webBridge;
};

export const nativeBridge = getNativeBridge();

declare global {
  interface Window {
    lynxBridge?: NativeBridge;
  }
}
