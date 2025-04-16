import React from "react";
import { useSharedState } from "./shared-state-provider";
import { cn } from "../../utils/cn";
import { AnimatedTabs } from "./aceternity/animated-tabs";
import { SpotlightCard } from "./aceternity/spotlight-card";

export type AgentType = "control-plane" | "worker" | "ui" | "codegen" | "scaffolding";

interface AgentViewSwitcherProps {
  className?: string;
}

export const AgentViewSwitcher: React.FC<AgentViewSwitcherProps> = ({
  className,
}) => {
  const { state, updateState } = useSharedState();
  
  const agentTypes: { id: AgentType; label: string; description: string }[] = [
    {
      id: "control-plane",
      label: "Control Plane Agent",
      description: "Manages and orchestrates all worker agents",
    },
    {
      id: "worker",
      label: "Worker Agent",
      description: "Executes specific tasks assigned by the control plane",
    },
    {
      id: "ui",
      label: "UI Agent",
      description: "Specializes in UI/UX development tasks",
    },
    {
      id: "codegen",
      label: "Codegen Agent",
      description: "Generates code based on specifications",
    },
    {
      id: "scaffolding",
      label: "Scaffolding Agent",
      description: "Creates application structure and boilerplate",
    },
  ];

  const activeAgentType = state.activeAgentId 
    ? (agentTypes.find(a => a.id === state.activeAgentId)?.id || "control-plane") 
    : "control-plane";

  const handleAgentTypeChange = (agentType: AgentType) => {
    updateState({ activeAgentId: agentType });
  };

  const tabs = agentTypes.map(agent => ({
    id: agent.id,
    label: agent.label,
    content: (
      <SpotlightCard className="p-4">
        <h3 className="text-lg font-medium mb-2">{agent.label}</h3>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">{agent.description}</p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="border rounded-md p-3 bg-white/5">
            <h4 className="text-sm font-medium mb-2">Status</h4>
            <div className="flex items-center gap-2">
              <span className={cn(
                "h-2 w-2 rounded-full",
                agent.id === activeAgentType ? "bg-emerald-500" : "bg-gray-400"
              )}></span>
              <span className="text-sm">
                {agent.id === activeAgentType ? "Active" : "Inactive"}
              </span>
            </div>
          </div>
          
          <div className="border rounded-md p-3 bg-white/5">
            <h4 className="text-sm font-medium mb-2">Resources</h4>
            <div className="text-sm">
              <div className="flex justify-between">
                <span>CPU:</span>
                <span>25%</span>
              </div>
              <div className="flex justify-between">
                <span>Memory:</span>
                <span>128MB</span>
              </div>
            </div>
          </div>
        </div>
      </SpotlightCard>
    ),
  }));

  return (
    <div className={cn("border rounded-lg overflow-hidden", className)}>
      <div className="bg-gray-50 dark:bg-gray-800 p-4 border-b">
        <h2 className="text-lg font-medium">Agent View Switcher</h2>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Toggle between different agent views
        </p>
      </div>
      
      <div className="p-4">
        <AnimatedTabs 
          tabs={tabs} 
          defaultTabId={activeAgentType}
        />
      </div>
    </div>
  );
};
