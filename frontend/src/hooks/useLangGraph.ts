import { useState, useEffect, useCallback } from "react";
import langGraphService, { 
  GraphState, 
  AgentConfig 
} from "../services/langgraph-service";
import { useSharedState } from "../components/ui/shared-state-provider";

export interface UseLangGraphOptions {
  autoConnect?: boolean;
  initialAgentType?: AgentConfig["agentType"];
  initialModelId?: AgentConfig["modelId"];
}

export function useLangGraph(options: UseLangGraphOptions = {}) {
  const { 
    autoConnect = false, 
    initialAgentType = "swe-agent",
    initialModelId = "gemini-2.5-pro"
  } = options;
  
  const [graphId, setGraphId] = useState<string | null>(null);
  const [graphState, setGraphState] = useState<GraphState | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const { state: sharedState, updateState } = useSharedState();

  const createGraph = useCallback(async (config: AgentConfig) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const state = await langGraphService.createGraph(config);
      setGraphId(state.graphId);
      setGraphState(state);
      
      updateState({
        activeAgentId: config.agentType,
        activeModelId: config.modelId,
        isGenerating: false,
        lastAction: "Graph created",
      });
      
      return state;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create graph";
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  }, [updateState]);

  const executeStep = useCallback(async (input: Record<string, any>) => {
    if (!graphId) {
      setError("No active graph. Create a graph first.");
      return null;
    }
    
    setIsLoading(true);
    setError(null);
    
    try {
      updateState({
        isGenerating: true,
        lastAction: "Executing graph step",
      });
      
      const result = await langGraphService.executeGraphStep(graphId, input);
      setGraphState(result.state);
      
      updateState({
        isGenerating: false,
        lastAction: "Step executed",
        progress: 100,
      });
      
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to execute graph step";
      setError(errorMessage);
      
      updateState({
        isGenerating: false,
        lastError: errorMessage,
      });
      
      return null;
    } finally {
      setIsLoading(false);
    }
  }, [graphId, updateState]);

  useEffect(() => {
    if (!graphId) return;
    
    const unsubscribe = langGraphService.subscribeToGraphUpdates(
      graphId,
      (state) => {
        setGraphState(state);
        
        updateState({
          isGenerating: state.status === "running",
          progress: state.status === "completed" ? 100 : 50,
          lastAction: `Graph state updated: ${state.status}`,
          lastError: state.error || null,
        });
      },
      (error) => {
        setError(error.message);
        updateState({
          lastError: error.message,
        });
      }
    );
    
    return unsubscribe;
  }, [graphId, updateState]);

  useEffect(() => {
    if (autoConnect && !graphId) {
      createGraph({
        agentType: initialAgentType,
        modelId: initialModelId,
      });
    }
  }, [autoConnect, createGraph, graphId, initialAgentType, initialModelId]);

  const switchAgent = useCallback(async (agentType: AgentConfig["agentType"]) => {
    if (graphId && graphState) {
      return createGraph({
        agentType,
        modelId: sharedState.activeModelId as AgentConfig["modelId"] || "gemini-2.5-pro",
        initialState: graphState.state,
      });
    }
    return null;
  }, [graphId, graphState, createGraph, sharedState.activeModelId]);

  const switchModel = useCallback(async (modelId: AgentConfig["modelId"]) => {
    if (graphId && graphState) {
      return createGraph({
        agentType: sharedState.activeAgentId as AgentConfig["agentType"] || "swe-agent",
        modelId,
        initialState: graphState.state,
      });
    }
    return null;
  }, [graphId, graphState, createGraph, sharedState.activeAgentId]);

  return {
    graphId,
    graphState,
    isLoading,
    error,
    createGraph,
    executeStep,
    switchAgent,
    switchModel,
  };
}

export default useLangGraph;
