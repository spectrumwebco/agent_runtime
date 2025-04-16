import React, { useState, useEffect } from "react";
import { cn } from "../../../utils/cn";
import { GradientCard } from "../aceternity/gradient-card";
import { TextGradient } from "../aceternity/text-gradient";

interface Task {
  id: string;
  name: string;
  status: "pending" | "in-progress" | "completed" | "failed";
  progress: number;
  startTime?: Date;
  endTime?: Date;
}

interface AITaskTrackerProps {
  className?: string;
  tasks: Task[];
  currentTaskId?: string;
  onTaskClick?: (taskId: string) => void;
}

export const AITaskTracker: React.FC<AITaskTrackerProps> = ({
  className,
  tasks,
  currentTaskId,
  onTaskClick,
}) => {
  const [expandedTaskId, setExpandedTaskId] = useState<string | null>(null);
  
  useEffect(() => {
    if (currentTaskId) {
      setExpandedTaskId(currentTaskId);
    }
  }, [currentTaskId]);

  const toggleTask = (taskId: string) => {
    setExpandedTaskId(expandedTaskId === taskId ? null : taskId);
    onTaskClick?.(taskId);
  };

  const getStatusColor = (status: Task["status"]) => {
    switch (status) {
      case "pending":
        return "text-gray-500";
      case "in-progress":
        return "text-blue-500";
      case "completed":
        return "text-emerald-500";
      case "failed":
        return "text-red-500";
      default:
        return "text-gray-500";
    }
  };

  const getStatusIcon = (status: Task["status"]) => {
    switch (status) {
      case "pending":
        return "⏳";
      case "in-progress":
        return "🔄";
      case "completed":
        return "✅";
      case "failed":
        return "❌";
      default:
        return "⏳";
    }
  };

  const formatTime = (date?: Date) => {
    if (!date) return "N/A";
    return date.toLocaleTimeString();
  };

  const calculateDuration = (start?: Date, end?: Date) => {
    if (!start) return "N/A";
    const endTime = end || new Date();
    const durationMs = endTime.getTime() - start.getTime();
    const seconds = Math.floor(durationMs / 1000);
    
    if (seconds < 60) {
      return `${seconds}s`;
    } else if (seconds < 3600) {
      return `${Math.floor(seconds / 60)}m ${seconds % 60}s`;
    } else {
      return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`;
    }
  };

  return (
    <GradientCard className={cn("p-4", className)}>
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold">
          <TextGradient>Kled Agent Tasks</TextGradient>
        </h3>
        <div className="text-sm text-gray-500">
          {tasks.filter(t => t.status === "completed").length}/{tasks.length} completed
        </div>
      </div>

      <div className="space-y-3">
        {tasks.map((task) => (
          <div 
            key={task.id}
            className={cn(
              "border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden transition-all duration-300",
              expandedTaskId === task.id ? "shadow-md" : "",
              "cursor-pointer"
            )}
            onClick={() => toggleTask(task.id)}
          >
            <div className="flex items-center justify-between p-3">
              <div className="flex items-center space-x-2">
                <span className="text-lg" aria-hidden="true">
                  {getStatusIcon(task.status)}
                </span>
                <span className={cn("font-medium", getStatusColor(task.status))}>
                  {task.name}
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-20 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-emerald-500 h-2 rounded-full transition-all duration-500 ease-in-out"
                    style={{ width: `${task.progress}%` }}
                  />
                </div>
                <span className="text-xs text-gray-500 dark:text-gray-400">
                  {task.progress}%
                </span>
              </div>
            </div>
            
            {expandedTaskId === task.id && (
              <div className="p-3 bg-gray-50 dark:bg-gray-800/50 border-t border-gray-200 dark:border-gray-700">
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="text-gray-500 dark:text-gray-400">Status:</div>
                  <div className={getStatusColor(task.status)}>
                    {task.status.charAt(0).toUpperCase() + task.status.slice(1)}
                  </div>
                  
                  <div className="text-gray-500 dark:text-gray-400">Start Time:</div>
                  <div>{formatTime(task.startTime)}</div>
                  
                  <div className="text-gray-500 dark:text-gray-400">End Time:</div>
                  <div>{formatTime(task.endTime)}</div>
                  
                  <div className="text-gray-500 dark:text-gray-400">Duration:</div>
                  <div>{calculateDuration(task.startTime, task.endTime)}</div>
                </div>
              </div>
            )}
          </div>
        ))}
        
        {tasks.length === 0 && (
          <div className="text-center py-6 text-gray-500 dark:text-gray-400">
            No tasks available
          </div>
        )}
      </div>
    </GradientCard>
  );
};

export default AITaskTracker;
