import React, { useState, useEffect } from "react";
import { useCompletion } from "@vercel/ai";
import { cn } from "../../../utils/cn";
import { GradientCard } from "../aceternity/gradient-card";
import { SpotlightButton } from "../aceternity/spotlight-button";

interface AIAgentStreamProps {
  className?: string;
  initialPrompt?: string;
  modelId?: string;
  onTaskComplete?: (result: string) => void;
  onTaskProgress?: (progress: number) => void;
}

export const AIAgentStream: React.FC<AIAgentStreamProps> = ({
  className,
  initialPrompt = "",
  modelId = "gemini-2.5-pro",
  onTaskComplete,
  onTaskProgress,
}) => {
  const [prompt, setPrompt] = useState(initialPrompt);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<
    "idle" | "running" | "completed" | "error"
  >("idle");

  const { completion, complete, isLoading, stop } = useCompletion({
    api: "/api/agent",
    body: {
      model: modelId,
    },
    onResponse: () => {
      setStatus("running");
    },
    onFinish: (result) => {
      setStatus("completed");
      setProgress(100);
      onTaskComplete?.(result);
    },
    onError: () => {
      setStatus("error");
    },
  });

  useEffect(() => {
    if (isLoading && completion) {
      const estimatedProgress = Math.min(
        Math.floor((completion.length / 500) * 100),
        99,
      );
      setProgress(estimatedProgress);
      onTaskProgress?.(estimatedProgress);
    }
  }, [completion, isLoading, onTaskProgress]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setProgress(0);
    setStatus("running");
    complete(prompt);
  };

  const handleStop = () => {
    stop();
    setStatus("idle");
  };

  return (
    <GradientCard className={cn("p-0 overflow-hidden", className)}>
      <div className="p-6">
        <h3 className="text-lg font-semibold mb-4">Kled.io Agent</h3>

        <form onSubmit={handleSubmit} className="mb-4">
          <textarea
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Describe the task for the agent..."
            className="w-full rounded-md border border-gray-300 dark:border-gray-600 bg-transparent p-3 focus:outline-none focus:ring-2 focus:ring-emerald-500 min-h-[100px]"
            disabled={isLoading}
          />

          <div className="flex justify-between mt-2">
            <div className="flex items-center">
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 mr-2 flex-grow max-w-[200px]">
                <div
                  className="bg-emerald-500 h-2 rounded-full transition-all duration-500 ease-in-out"
                  style={{ width: `${progress}%` }}
                />
              </div>
              <span className="text-xs text-gray-500 dark:text-gray-400">
                {progress}%
              </span>
            </div>

            {isLoading ? (
              <SpotlightButton onClick={handleStop} variant="outline" size="sm">
                Stop
              </SpotlightButton>
            ) : (
              <SpotlightButton
                type="submit"
                disabled={!prompt.trim() || isLoading}
              >
                Run Agent
              </SpotlightButton>
            )}
          </div>
        </form>
      </div>

      <div className="border-t border-gray-200 dark:border-gray-700 p-6 bg-gray-50 dark:bg-gray-800/50 min-h-[200px] max-h-[400px] overflow-y-auto">
        <div className="font-mono text-sm whitespace-pre-wrap">
          {completion || (
            <span className="text-gray-400 dark:text-gray-500">
              {status === "idle"
                ? "Agent is ready. Enter a task to begin."
                : status === "running"
                  ? "Agent is working on your task..."
                  : status === "error"
                    ? "An error occurred. Please try again."
                    : "Task completed."}
            </span>
          )}
        </div>
      </div>
    </GradientCard>
  );
};

export default AIAgentStream;
