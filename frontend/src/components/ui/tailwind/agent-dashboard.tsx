import React from "react";
import { cn } from "../../../utils/cn";
import { AgentCard } from "./agent-card";
import { AnimatedBackground } from "../aceternity/animated-background";

interface AgentInfo {
  id: string;
  title: string;
  description?: string;
  status: "idle" | "running" | "completed" | "error";
  progress: number;
  icon?: React.ReactNode;
}

interface AgentDashboardProps {
  className?: string;
  agents: AgentInfo[];
  onAgentClick?: (agentId: string) => void;
}

export const AgentDashboard: React.FC<AgentDashboardProps> = ({
  className,
  agents,
  onAgentClick,
}) => {
  return (
    <AnimatedBackground variant="dots" color="emerald">
      <div className={cn("p-6", className)}>
        <h2 className="text-2xl font-bold mb-6">Agent Dashboard</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {agents.map((agent) => (
            <AgentCard
              key={agent.id}
              title={agent.title}
              description={agent.description}
              status={agent.status}
              progress={agent.progress}
              icon={agent.icon}
              onClick={() => onAgentClick?.(agent.id)}
            />
          ))}
        </div>
        
        {agents.length === 0 && (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            <p>No agents found. Create a new agent to get started.</p>
          </div>
        )}
      </div>
    </AnimatedBackground>
  );
};

export default AgentDashboard;
