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
