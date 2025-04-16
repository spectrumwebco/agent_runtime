export type KledioEventType =
  | "message"
  | "agent_state_changed"
  | "run"
  | "read"
  | "write"
  | "edit"
  | "run_ipython"
  | "delegate"
  | "browse"
  | "browse_interactive"
  | "reject"
  | "think"
  | "finish"
  | "error"
  | "recall";

interface KledioBaseEvent {
  id: number;
  source: "agent" | "user";
  message: string;
  timestamp: string; // ISO 8601
}

export interface KledioActionEvent<T extends KledioEventType>
  extends KledioBaseEvent {
  action: T;
  args: Record<string, unknown>;
}

export interface KledioObservationEvent<T extends KledioEventType>
  extends KledioBaseEvent {
  cause: number;
  observation: T;
  content: string;
  extras: Record<string, unknown>;
}
