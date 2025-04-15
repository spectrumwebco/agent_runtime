
export interface FeedItem {
  type: string;
  message: string;
  format?: string;
  step?: number | null;
}

export interface RunConfig {
  environment: {
    data_path: string;
    image_name: string;
    script: string;
    repo_path: string;
    base_commit: string;
  };
  agent: {
    model: {
      model_name: string;
    };
  };
  extra: {
    test_run: boolean;
  };
}

export interface SharedStateEvent {
  type: string;
  state_type: string;
  state_id: string;
  data: any;
}

export interface SharedStateClient {
  connect: () => Promise<void>;
  subscribeToState: (stateType: string, stateId: string, callback: (state: any) => void) => void;
  unsubscribeFromState: (stateType: string, stateId: string) => void;
  sendEvent: (event: SharedStateEvent) => void;
}
