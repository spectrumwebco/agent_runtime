import React from "react";
import { useSharedState } from "./shared-state-provider";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";

interface StateVisualizationProps {
  className?: string;
}

export const StateVisualization: React.FC<StateVisualizationProps> = ({
  className,
}) => {
  const { state, updateState, resetState } = useSharedState();
  
  const stateHistory = [
    { timestamp: new Date(Date.now() - 3600000).toISOString(), label: "Initial State" },
    { timestamp: new Date(Date.now() - 2400000).toISOString(), label: "Tool Execution" },
    { timestamp: new Date(Date.now() - 1200000).toISOString(), label: "Agent Response" },
    { timestamp: new Date(Date.now()).toISOString(), label: "Current State" },
  ];

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <SpotlightCard className={cn("p-4", className)}>
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-medium">State Management</h3>
        <GradientButton 
          variant="outline" 
          size="sm"
          onClick={() => resetState()}
        >
          Reset State
        </GradientButton>
      </div>

      <div className="space-y-4">
        {/* Current State */}
        <div className="border rounded-md p-3 bg-white/5">
          <h4 className="text-sm font-medium mb-2">Current State</h4>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <div>Active Agent:</div>
            <div className="font-medium">{state.activeAgentId || "None"}</div>
            
            <div>Active Model:</div>
            <div className="font-medium">{state.activeModelId || "None"}</div>
            
            <div>Active View:</div>
            <div className="font-medium">{state.activeView}</div>
            
            <div>Is Generating:</div>
            <div className="font-medium">{state.isGenerating ? "Yes" : "No"}</div>
          </div>
        </div>

        {/* State History Timeline */}
        <div className="border rounded-md p-3 bg-white/5">
          <h4 className="text-sm font-medium mb-2">State History</h4>
          <div className="relative">
            {/* Timeline line */}
            <div className="absolute left-2.5 top-0 bottom-0 w-0.5 bg-gray-200 dark:bg-gray-700"></div>
            
            {/* Timeline points */}
            <div className="space-y-4">
              {stateHistory.map((item, index) => (
                <div key={index} className="flex items-start ml-2">
                  <div className={cn(
                    "h-5 w-5 rounded-full border-2 border-emerald-500 flex-shrink-0 z-10",
                    index === stateHistory.length - 1 ? "bg-emerald-500" : "bg-background"
                  )}></div>
                  <div className="ml-3">
                    <div className="text-sm font-medium">{item.label}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {formatTimestamp(item.timestamp)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* State Database Info */}
        <div className="border rounded-md p-3 bg-white/5">
          <h4 className="text-sm font-medium mb-2">State Database</h4>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <div>Main DB:</div>
            <div className="font-medium text-emerald-500">Connected</div>
            
            <div>Read-only DB:</div>
            <div className="font-medium text-emerald-500">Connected</div>
            
            <div>Rollback DB:</div>
            <div className="font-medium text-emerald-500">Connected</div>
          </div>
        </div>
      </div>
    </SpotlightCard>
  );
};
