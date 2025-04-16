import React from "react";
import { cn } from "../../../utils/cn";
import { GradientCard } from "../../ui/aceternity/gradient-card";

interface AgentCardProps {
  className?: string;
  title: string;
  description?: string;
  status: "idle" | "running" | "completed" | "error";
  progress: number;
  icon?: React.ReactNode;
  onClick?: () => void;
}

export const AgentCard: React.FC<AgentCardProps> = ({
  className,
  title,
  description,
  status,
  progress,
  icon,
  onClick,
}) => {
  const statusColors = {
    idle: "bg-gray-200 dark:bg-gray-700",
    running: "bg-blue-200 dark:bg-blue-700",
    completed: "bg-emerald-200 dark:bg-emerald-700",
    error: "bg-red-200 dark:bg-red-700",
  };

  const statusText = {
    idle: "Idle",
    running: "Running",
    completed: "Completed",
    error: "Error",
  };

  return (
    <GradientCard className={cn("cursor-pointer", className)} onClick={onClick}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <h3 className="text-lg font-semibold">{title}</h3>
            <span
              className={cn(
                "text-xs px-2 py-0.5 rounded-full",
                statusColors[status],
              )}
            >
              {statusText[status]}
            </span>
          </div>
          {description && (
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {description}
            </p>
          )}
        </div>
        {icon && <div className="text-2xl">{icon}</div>}
      </div>

      <div className="mt-4">
        <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mb-1">
          <span>Progress</span>
          <span>{progress}%</span>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div
            className="bg-emerald-500 h-2 rounded-full transition-all duration-500 ease-in-out"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>
    </GradientCard>
  );
};

export default AgentCard;
