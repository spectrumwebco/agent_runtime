import React from "react";
import { cn } from "../../utils/cn";
import { useSharedState } from "./shared-state-provider";

interface StateVisualizationProps {
  className?: string;
}

export const StateVisualization: React.FC<StateVisualizationProps> = ({
  className,
}) => {
  const { state } = useSharedState();

  const systemResources = {
    cpu: 45,
    memory: 32,
    network: 18,
  };

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 shadow-md p-4",
        className,
      )}
    >
      <h3 className="text-lg font-medium mb-2">State Visualization</h3>
      <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
        Current state of the Kled.io agent system
      </p>

      <div className="space-y-4">
        <div className="border rounded-md p-3 bg-gray-50 dark:bg-gray-900">
          <h4 className="text-sm font-medium mb-2">Agent State</h4>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <div>Status:</div>
            <div className="font-mono">
              {state.isGenerating ? "Active" : "Inactive"}
            </div>

            <div>View Mode:</div>
            <div className="font-mono">
              {state.viewMode === "control" ? "Control Plane" : "Worker Agents"}
            </div>

            <div>Active View:</div>
            <div className="font-mono">{state.activeView}</div>

            <div>Selected Model:</div>
            <div className="font-mono">{state.activeModelId || "None"}</div>

            <div>Progress:</div>
            <div className="font-mono">{state.progress}%</div>

            <div>Last Action:</div>
            <div className="font-mono">{state.lastAction || "None"}</div>
          </div>
        </div>

        <div className="border rounded-md p-3 bg-gray-50 dark:bg-gray-900">
          <h4 className="text-sm font-medium mb-2">Active Tools</h4>
          {state.activeToolIds.length > 0 ? (
            <ul className="space-y-1 text-sm">
              {state.activeToolIds.map((toolId) => (
                <li key={toolId} className="border-l-2 border-emerald-500 pl-2">
                  {toolId}
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-gray-500 dark:text-gray-400">
              No active tools
            </p>
          )}
        </div>

        <div className="border rounded-md p-3 bg-gray-50 dark:bg-gray-900">
          <h4 className="text-sm font-medium mb-2">System Resources</h4>
          <div className="space-y-2">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span>CPU Usage</span>
                <span>{systemResources.cpu}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-emerald-500 h-2 rounded-full"
                  style={{ width: `${systemResources.cpu}%` }}
                ></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-1">
                <span>Memory Usage</span>
                <span>{systemResources.memory}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-blue-500 h-2 rounded-full"
                  style={{ width: `${systemResources.memory}%` }}
                ></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-1">
                <span>Network Usage</span>
                <span>{systemResources.network}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-purple-500 h-2 rounded-full"
                  style={{ width: `${systemResources.network}%` }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StateVisualization;
