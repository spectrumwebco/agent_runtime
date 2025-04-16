import React from "react";
import { cn } from "../../utils/cn";
import { useSharedState } from "./shared-state-provider";

interface AgentViewSwitcherProps {
  className?: string;
}

export const AgentViewSwitcher: React.FC<AgentViewSwitcherProps> = ({
  className,
}) => {
  const { state, updateState } = useSharedState();

  const viewModes = [
    { id: "control", label: "Control Plane" },
    { id: "worker", label: "Worker Agents" },
  ];

  const currentViewMode = state.viewMode || "control";

  const handleViewModeChange = (viewMode: string) => {
    updateState({
      activeView: viewMode === "control" ? "code" : "terminal",
      viewMode: viewMode as "control" | "worker",
    });
  };

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 shadow-md p-4",
        className,
      )}
    >
      <div className="flex flex-col space-y-4">
        <div>
          <h3 className="text-lg font-medium mb-1">Agent View Switcher</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Toggle between Control Plane Agent and Worker Agents
          </p>
        </div>

        <div className="flex gap-3">
          {viewModes.map((mode) => (
            <button
              key={mode.id}
              onClick={() => handleViewModeChange(mode.id)}
              className={cn(
                "relative inline-flex items-center justify-center px-4 py-2 rounded-md font-medium text-sm transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500",
                currentViewMode === mode.id
                  ? "bg-gradient-to-r from-emerald-400 to-emerald-600 text-white border border-emerald-500 shadow-md"
                  : "bg-transparent border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800",
                "flex-1",
              )}
            >
              <span className="relative z-10">{mode.label}</span>
            </button>
          ))}
        </div>

        <div className="text-sm text-gray-500 dark:text-gray-400 pt-2 border-t border-gray-200 dark:border-gray-700">
          <p>
            {currentViewMode === "control"
              ? "Control Plane Agent manages task coordination and delegation"
              : "Worker Agents execute specialized tasks and report back to Control Plane"}
          </p>
        </div>
      </div>
    </div>
  );
};

export default AgentViewSwitcher;
