import React from "react";
import { cn } from "../../utils/cn";
import { KledHeader } from "./kled-header";
import { FlippedLayout } from "./flipped-layout";
import { MultiAgentSelector } from "./multi-agent-selector";
import { AIModelSelector } from "./ai-model-selector";
import {
  LangGraphVisualizer,
  GraphNode,
  GraphEdge,
} from "./langgraph-visualizer";
import { useSharedState } from "./shared-state-provider";

interface KledMainLayoutProps {
  className?: string;
  children?: React.ReactNode;
}

export const KledMainLayout: React.FC<KledMainLayoutProps> = ({
  className,
  children,
}) => {
  const agents = [
    {
      id: "swe-agent",
      name: "Software Engineering Agent",
      type: "software" as const,
      description: "Specialized in code development and engineering tasks",
    },
    {
      id: "scaffolding-agent",
      name: "App Scaffolding Agent",
      type: "scaffolding" as const,
      description: "Creates application structure and boilerplate",
    },
    {
      id: "ui-agent",
      name: "UI/UX Agent",
      type: "ui" as const,
      description: "Designs and implements user interfaces",
    },
    {
      id: "codegen-agent",
      name: "Codegen Agent",
      type: "codegen" as const,
      description: "Generates code based on specifications",
    },
  ];

  const models = [
    {
      id: "gemini-2.5-pro",
      name: "Gemini 2.5 Pro",
      provider: "Google",
      description: "Advanced model for coding tasks with 1M context window",
      isVerified: true,
      capabilities: ["Coding", "Reasoning", "Planning"],
    },
    {
      id: "llama-4",
      name: "Llama 4",
      provider: "Meta",
      description: "Specialized for operations and reasoning tasks",
      isVerified: true,
      capabilities: ["Operations", "Reasoning"],
    },
    {
      id: "openai-gpt4",
      name: "GPT-4",
      provider: "OpenAI",
      description: "General purpose AI model",
      isVerified: true,
      capabilities: ["General", "Coding", "Reasoning"],
    },
  ];

  const nodes: GraphNode[] = [
    { id: "swe", type: "agent", label: "SWE Agent", status: "active" },
    { id: "ui", type: "agent", label: "UI Agent", status: "idle" },
    { id: "codegen", type: "agent", label: "Codegen Agent", status: "idle" },
    {
      id: "scaffolding",
      type: "agent",
      label: "Scaffolding Agent",
      status: "idle",
    },
    {
      id: "code-tool",
      type: "tool",
      label: "Code Analysis",
      status: "completed",
    },
    { id: "search-tool", type: "tool", label: "Search", status: "idle" },
  ];

  const edges: GraphEdge[] = [
    { source: "swe", target: "code-tool", label: "uses" },
    { source: "swe", target: "ui", label: "delegates to" },
    { source: "ui", target: "codegen", label: "requests from" },
    { source: "codegen", target: "scaffolding", label: "builds with" },
  ];

  const { state, updateState } = useSharedState();

  const controlsContent = (
    <div className="p-4 space-y-6">
      <h2 className="text-xl font-bold mb-4">Kled Control Panel</h2>

      <MultiAgentSelector
        agents={agents}
        selectedAgentId={state.activeAgentId || "swe-agent"}
        onSelectAgent={(agentId) => updateState({ activeAgentId: agentId })}
      />

      <AIModelSelector
        models={models}
        selectedModelId={state.activeModelId || "gemini-2.5-pro"}
        onSelectModel={(modelId) => updateState({ activeModelId: modelId })}
      />

      <LangGraphVisualizer nodes={nodes} edges={edges} />
    </div>
  );

  return (
    <div className={cn("flex flex-col h-screen", className)}>
      <KledHeader />

      <div className="flex-1 overflow-hidden">
        <FlippedLayout
          editor={
            <div className="h-full w-full bg-gray-50 dark:bg-gray-900 p-4">
              {children || (
                <div className="flex items-center justify-center h-full">
                  <p className="text-gray-500">IDE/Code Editor Area</p>
                </div>
              )}
            </div>
          }
          controls={controlsContent}
        />
      </div>
    </div>
  );
};

export default KledMainLayout;
