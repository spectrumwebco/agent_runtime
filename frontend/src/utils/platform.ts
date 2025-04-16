/**
 * Platform detection and utilities for cross-platform support
 * Enables consistent behavior across React web, Electron desktop, and Lynx-React mobile
 */

export type Platform = "web" | "electron" | "mobile";

/**
 * Detect the current platform
 * @returns The current platform: 'web', 'electron', or 'mobile'
 */
export const detectPlatform = (): Platform => {
  if (typeof window !== "undefined" && window.process && window.process.type) {
    return "electron";
  }

  if (
    typeof navigator !== "undefined" &&
    (navigator.userAgent.includes("LynxReact") ||
      (typeof window !== "undefined" && window.lynxReact))
  ) {
    return "mobile";
  }

  return "web";
};

/**
 * Current platform
 */
export const currentPlatform = detectPlatform();

/**
 * Check if running on web platform
 */
export const isWeb = currentPlatform === "web";

/**
 * Check if running on Electron desktop platform
 */
export const isElectron = currentPlatform === "electron";

/**
 * Check if running on Lynx-React mobile platform
 */
export const isMobile = currentPlatform === "mobile";

/**
 * Platform-specific feature detection
 */
export const platformFeatures = {
  common: {
    supportsLocalStorage:
      typeof window !== "undefined" && !!window.localStorage,
    supportsServiceWorker:
      typeof navigator !== "undefined" && "serviceWorker" in navigator,
  },

  web: {
    supportsSharing: typeof navigator !== "undefined" && !!navigator.share,
  },

  electron: {
    supportsFileSystem: isElectron,
    supportsNativeMenus: isElectron,
  },

  mobile: {
    supportsTouch:
      typeof navigator !== "undefined" &&
      "maxTouchPoints" in navigator &&
      navigator.maxTouchPoints > 0,
    supportsGeolocation:
      typeof navigator !== "undefined" && !!navigator.geolocation,
    supportsCamera:
      typeof navigator !== "undefined" &&
      !!navigator.mediaDevices &&
      !!navigator.mediaDevices.getUserMedia,
  },
};

/**
 * Run platform-specific code
 * @param options Object containing platform-specific implementations
 * @returns Result of the platform-specific implementation
 */
export function runForPlatform<T>(options: {
  web?: () => T;
  electron?: () => T;
  mobile?: () => T;
  fallback?: () => T;
}): T | undefined {
  const { web, electron, mobile, fallback } = options;

  if (isWeb && web) {
    return web();
  } else if (isElectron && electron) {
    return electron();
  } else if (isMobile && mobile) {
    return mobile();
  } else if (fallback) {
    return fallback();
  }

  return undefined;
}
