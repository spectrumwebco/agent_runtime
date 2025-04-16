import React, { useState, useEffect } from "react";
import { cn } from "../../../utils/cn";
import { GradientCard } from "../aceternity/gradient-card";
import { AnimatedBackground } from "../aceternity/animated-background";

interface StateNode {
  id: string;
  name: string;
  value: any;
  children?: StateNode[];
}

interface AIStateVisualizerProps {
  className?: string;
  state: StateNode[];
  title?: string;
  onStateNodeClick?: (nodeId: string) => void;
}

export const AIStateVisualizer: React.FC<AIStateVisualizerProps> = ({
  className,
  state,
  title = "Kled Agent State",
  onStateNodeClick,
}) => {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());

  const toggleNode = (nodeId: string) => {
    const newExpandedNodes = new Set(expandedNodes);
    if (newExpandedNodes.has(nodeId)) {
      newExpandedNodes.delete(nodeId);
    } else {
      newExpandedNodes.add(nodeId);
    }
    setExpandedNodes(newExpandedNodes);
    onStateNodeClick?.(nodeId);
  };

  const renderStateNode = (node: StateNode, depth = 0) => {
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.id} className="mb-1">
        <div
          className={cn(
            "flex items-start py-1 px-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800 cursor-pointer",
            "transition-colors duration-200",
          )}
          style={{ paddingLeft: `${depth * 16 + 8}px` }}
          onClick={() => toggleNode(node.id)}
        >
          {hasChildren && (
            <span
              className="mr-2 text-gray-500 dark:text-gray-400 transform transition-transform duration-200"
              style={{
                transform: isExpanded ? "rotate(90deg)" : "rotate(0deg)",
              }}
            >
              ▶
            </span>
          )}
          <div className="flex-1">
            <span className="font-medium text-emerald-600 dark:text-emerald-400">
              {node.name}
            </span>
            {!hasChildren && (
              <span className="ml-2 text-gray-600 dark:text-gray-300">
                {typeof node.value === "object"
                  ? JSON.stringify(node.value).substring(0, 50) +
                    (JSON.stringify(node.value).length > 50 ? "..." : "")
                  : String(node.value)}
              </span>
            )}
          </div>
        </div>

        {isExpanded && hasChildren && (
          <div className="ml-4">
            {node.children!.map((child) => renderStateNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <GradientCard className={cn("overflow-hidden", className)}>
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-lg font-semibold">{title}</h3>
      </div>

      <AnimatedBackground variant="dots" color="emerald" animate={false}>
        <div className="p-4 max-h-[400px] overflow-y-auto">
          {state.length > 0 ? (
            <div className="space-y-2">
              {state.map((node) => renderStateNode(node))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
              No state data available
            </div>
          )}
        </div>
      </AnimatedBackground>
    </GradientCard>
  );
};

export default AIStateVisualizer;
