import { GraphNode, GraphEdge } from "../components/ui/langgraph-visualizer";

export interface GraphState {
  graphId: string;
  nodes: GraphNode[];
  edges: GraphEdge[];
  currentNode: string;
  state: Record<string, any>;
  status: "created" | "running" | "completed" | "error";
  error?: string;
}

export interface AgentConfig {
  agentType: "swe-agent" | "ui-agent" | "scaffolding-agent" | "codegen-agent";
  modelId: "gemini-2.5-pro" | "llama-4" | "openai-gpt4";
  initialState?: Record<string, any>;
}

class LangGraphService {
  private baseUrl: string;
  private websocketUrl: string;
  private activeWebsockets: Map<string, WebSocket> = new Map();

  constructor() {
    this.baseUrl = "/api/langgraph";
    this.websocketUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws/langgraph`;
  }

  /**
   * Create a new agent graph
   */
  async createGraph(config: AgentConfig): Promise<GraphState> {
    const response = await fetch(`${this.baseUrl}/create`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(config),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to create graph");
    }

    const data = await response.json();
    return data.state;
  }

  /**
   * Get the current state of a graph
   */
  async getGraphState(graphId: string): Promise<GraphState> {
    const response = await fetch(`${this.baseUrl}/state/${graphId}`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to get graph state");
    }

    return await response.json();
  }

  /**
   * Execute a step in the graph
   */
  async executeGraphStep(graphId: string, input: Record<string, any>): Promise<{
    result: Record<string, any>;
    state: GraphState;
  }> {
    const response = await fetch(`${this.baseUrl}/execute/${graphId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(input),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to execute graph step");
    }

    return await response.json();
  }

  /**
   * Subscribe to real-time updates for a graph
   */
  subscribeToGraphUpdates(
    graphId: string,
    onUpdate: (state: GraphState) => void,
    onError?: (error: Error) => void
  ): () => void {
    if (this.activeWebsockets.has(graphId)) {
      this.activeWebsockets.get(graphId)?.close();
    }

    const ws = new WebSocket(`${this.websocketUrl}/${graphId}`);

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onUpdate(data);
      } catch (error) {
        if (onError) {
          onError(new Error(`Failed to parse websocket message: ${error}`));
        }
      }
    };

    ws.onerror = (event) => {
      if (onError) {
        onError(new Error("WebSocket error"));
      }
    };

    this.activeWebsockets.set(graphId, ws);

    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
      this.activeWebsockets.delete(graphId);
    };
  }
}

export const langGraphService = new LangGraphService();
export default langGraphService;
