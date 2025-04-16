import React from "react";
import { cn } from "../../utils/cn";

export interface Agent {
  id: string;
  name: string;
  type: "software" | "scaffolding" | "ui" | "codegen";
  description: string;
  icon?: React.ReactNode;
}

interface MultiAgentSelectorProps {
  agents: Agent[];
  selectedAgentId: string;
  onSelectAgent: (agentId: string) => void;
  className?: string;
}

export const MultiAgentSelector: React.FC<MultiAgentSelectorProps> = ({
  agents,
  selectedAgentId,
  onSelectAgent,
  className,
}) => {
  return (
    <div className={cn("flex flex-col space-y-2", className)}>
      <h3 className="text-lg font-medium mb-2">Agent Selection</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
        {agents.map((agent) => (
          <button
            key={agent.id}
            onClick={() => onSelectAgent(agent.id)}
            className={cn(
              "aceternity-card p-4 text-left transition-all",
              selectedAgentId === agent.id
                ? "border-emerald-500 ring-2 ring-emerald-500"
                : "border-gray-200 dark:border-gray-700 hover:border-emerald-500"
            )}
          >
            <div className="flex items-center gap-3">
              {agent.icon && <div className="text-2xl">{agent.icon}</div>}
              <div>
                <h4 className="font-medium">{agent.name}</h4>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {agent.description}
                </p>
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};

export default MultiAgentSelector;
