import React from "react";

interface ProgressBarProps {
  progress: number;
  status?: "idle" | "running" | "paused" | "completed" | "failed";
  showPercentage?: boolean;
  className?: string;
  height?: number;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  status = "running",
  showPercentage = true,
  className = "",
  height = 8,
}) => {
  const normalizedProgress = Math.min(Math.max(progress, 0), 100);

  const statusColors = {
    idle: "bg-gray-400",
    running: "bg-emerald-500",
    paused: "bg-yellow-500",
    completed: "bg-emerald-500",
    failed: "bg-red-500",
  };

  const statusText = {
    idle: "Idle",
    running: "Running",
    paused: "Paused",
    completed: "Completed",
    failed: "Failed",
  };

  return (
    <div className={`w-full ${className}`}>
      <div className="flex justify-between items-center mb-1">
        {status && (
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {statusText[status]}
          </span>
        )}
        {showPercentage && (
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {normalizedProgress.toFixed(0)}%
          </span>
        )}
      </div>
      <div
        className="w-full bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden"
        style={{ height: `${height}px` }}
      >
        <div
          className={`${statusColors[status]} transition-all duration-300 ease-in-out`}
          style={{ width: `${normalizedProgress}%`, height: "100%" }}
        ></div>
      </div>
    </div>
  );
};

export default ProgressBar;
