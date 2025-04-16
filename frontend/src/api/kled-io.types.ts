export interface Feedback {
  rating: number;
  comment?: string;
  tags?: string[];
}

export interface FeedbackResponse {
  id: string;
  rating: number;
  comment?: string;
  tags?: string[];
  created_at: string;
}

export interface GitHubAccessTokenResponse {
  access_token: string;
  token_type: string;
  scope: string;
}

export interface GetConfigResponse {
  APP_MODE: "saas" | "oss";
  FEATURE_FLAGS: {
    ENABLE_BILLING: boolean;
    ENABLE_ANALYTICS: boolean;
    ENABLE_SECURITY_ANALYZER: boolean;
    ENABLE_GITHUB_APP: boolean;
    ENABLE_VSCODE: boolean;
  };
  GITHUB_APP_URL?: string;
  GITHUB_CLIENT_ID?: string;
  GITHUB_REDIRECT_URI?: string;
}

export interface GetVSCodeUrlResponse {
  url: string;
}

export interface AuthenticateResponse {
  status: number;
  message: string;
}

export interface Conversation {
  conversation_id: string;
  title: string;
  created_at: string;
  updated_at: string;
  repository?: {
    id: number;
    name: string;
    full_name: string;
    private: boolean;
    html_url: string;
    description: string;
    fork: boolean;
    url: string;
    created_at: string;
    updated_at: string;
    pushed_at: string;
    git_url: string;
    ssh_url: string;
    clone_url: string;
    homepage: string;
    size: number;
    stargazers_count: number;
    watchers_count: number;
    language: string;
    forks_count: number;
    archived: boolean;
    disabled: boolean;
    open_issues_count: number;
    license: {
      key: string;
      name: string;
      url: string;
      spdx_id: string;
      node_id: string;
    };
    allow_forking: boolean;
    is_template: boolean;
    topics: string[];
    visibility: string;
    forks: number;
    open_issues: number;
    watchers: number;
    default_branch: string;
  };
}

export interface ResultSet<T> {
  count: number;
  next: string | null;
  previous: string | null;
  results: T[];
}

export interface GetTrajectoryResponse {
  nodes: {
    id: string;
    label: string;
    type: string;
  }[];
  edges: {
    from: string;
    to: string;
    label: string;
  }[];
}
