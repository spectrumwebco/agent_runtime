import React, { useEffect } from "react";
import { useLangGraph } from "../../hooks/useLangGraph";
import { useSharedState } from "./shared-state-provider";
import { LangGraphVisualizer } from "./langgraph-visualizer";
import { Button } from "./shadcn/button";

interface LangGraphIntegrationProps {
  className?: string;
}

export const LangGraphIntegration: React.FC<LangGraphIntegrationProps> = ({
  className,
}) => {
  const { state, updateState } = useSharedState();

  const {
    graphId,
    graphState,
    isLoading,
    error,
    createGraph,
    executeStep,
    switchAgent,
    switchModel,
  } = useLangGraph({
    autoConnect: true,
    initialAgentType: "swe-agent",
    initialModelId: "gemini-2.5-pro",
  });

  useEffect(() => {
    if (state.activeAgentId && graphId) {
      switchAgent(state.activeAgentId as any);
    }
  }, [state.activeAgentId, graphId, switchAgent]);

  useEffect(() => {
    if (state.activeModelId && graphId) {
      switchModel(state.activeModelId as any);
    }
  }, [state.activeModelId, graphId, switchModel]);

  const handleExecuteStep = async () => {
    if (!graphId) {
      await createGraph({
        agentType: (state.activeAgentId as any) || "swe-agent",
        modelId: (state.activeModelId as any) || "gemini-2.5-pro",
      });
    }

    await executeStep({
      input: "Execute next step in the agent workflow",
    });
  };

  return (
    <div className={className}>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-medium">Agent Graph</h3>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                createGraph({
                  agentType: (state.activeAgentId as any) || "swe-agent",
                  modelId: (state.activeModelId as any) || "gemini-2.5-pro",
                })
              }
              disabled={isLoading}
            >
              Reset Graph
            </Button>
            <Button
              variant="emerald"
              size="sm"
              onClick={handleExecuteStep}
              disabled={isLoading}
            >
              {isLoading ? "Processing..." : "Execute Step"}
            </Button>
          </div>
        </div>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 text-red-800 dark:text-red-300 p-3 rounded-md text-sm">
            {error}
          </div>
        )}

        {graphState ? (
          <LangGraphVisualizer
            nodes={graphState.nodes}
            edges={graphState.edges}
          />
        ) : (
          <div className="border rounded-lg p-8 flex items-center justify-center">
            <p className="text-gray-500 dark:text-gray-400">
              {isLoading
                ? "Creating agent graph..."
                : "No active agent graph. Click 'Reset Graph' to create one."}
            </p>
          </div>
        )}

        {graphState && (
          <div className="border rounded-lg p-4 bg-gray-50 dark:bg-gray-900">
            <h4 className="text-sm font-medium mb-2">Current State</h4>
            <pre className="text-xs overflow-auto max-h-40 bg-white dark:bg-gray-800 p-2 rounded">
              {JSON.stringify(graphState.state, null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};

export default LangGraphIntegration;
