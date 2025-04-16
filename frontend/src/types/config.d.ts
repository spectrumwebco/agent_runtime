interface Config {
  POSTHOG_CLIENT_KEY?: string;
  API_URL?: string;
  WEBSOCKET_URL?: string;
  ENVIRONMENT?: string;
  DEBUG_MODE?: boolean;
  [key: string]: any;
}

declare module '@/hooks/query/use-config' {
  export function useConfig(): { data?: Config; isLoading: boolean; error: any };
}
