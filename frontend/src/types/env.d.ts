interface ImportMetaEnv {
  readonly VITE_MOCK_API: string;
  readonly VITE_MOCK_SAAS: string;
  readonly APP_MODE: string;
  readonly APP_SLUG: string;
  readonly GITHUB_CLIENT_ID: string;
  readonly FEATURE_FLAGS: Record<string, boolean>;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare global {
  interface Window {
    APP_MODE: string;
    APP_SLUG: string;
    GITHUB_CLIENT_ID: string;
    FEATURE_FLAGS: Record<string, boolean>;
  }
  
  const APP_MODE: string;
  const APP_SLUG: string;
  const GITHUB_CLIENT_ID: string;
  const FEATURE_FLAGS: Record<string, boolean>;
}
