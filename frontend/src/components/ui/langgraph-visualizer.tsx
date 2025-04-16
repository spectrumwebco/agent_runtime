import React from "react";
import { cn } from "../../utils/cn";

export interface GraphNode {
  id: string;
  type: "agent" | "tool" | "state" | "action";
  label: string;
  status?: "idle" | "active" | "completed" | "error";
}

export interface GraphEdge {
  source: string;
  target: string;
  label?: string;
}

interface LangGraphVisualizerProps {
  nodes: GraphNode[];
  edges: GraphEdge[];
  className?: string;
}

export const LangGraphVisualizer: React.FC<LangGraphVisualizerProps> = ({
  nodes,
  edges,
  className,
}) => {
  return (
    <div className={cn("border rounded-lg p-4 bg-gray-50 dark:bg-gray-900", className)}>
      <h3 className="text-lg font-medium mb-4">Agent Graph Visualization</h3>
      
      <div className="grid grid-cols-1 gap-4">
        <div className="border rounded p-3 bg-white dark:bg-gray-800">
          <h4 className="font-medium mb-2">Nodes</h4>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
            {nodes.map((node) => (
              <div
                key={node.id}
                className={cn(
                  "p-2 rounded-md text-sm border",
                  node.status === "active" && "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20",
                  node.status === "completed" && "border-blue-500 bg-blue-50 dark:bg-blue-900/20",
                  node.status === "error" && "border-red-500 bg-red-50 dark:bg-red-900/20",
                  node.status === "idle" && "border-gray-200 dark:border-gray-700"
                )}
              >
                <div className="font-medium">{node.label}</div>
                <div className="text-xs text-gray-500 dark:text-gray-400">{node.type}</div>
              </div>
            ))}
          </div>
        </div>
        
        <div className="border rounded p-3 bg-white dark:bg-gray-800">
          <h4 className="font-medium mb-2">Connections</h4>
          <div className="space-y-1">
            {edges.map((edge, index) => (
              <div key={index} className="text-sm flex items-center gap-2">
                <span className="font-medium">{nodes.find(n => n.id === edge.source)?.label}</span>
                <span className="text-gray-500">→</span>
                <span className="font-medium">{nodes.find(n => n.id === edge.target)?.label}</span>
                {edge.label && (
                  <span className="text-xs text-gray-500 dark:text-gray-400">({edge.label})</span>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default LangGraphVisualizer;
