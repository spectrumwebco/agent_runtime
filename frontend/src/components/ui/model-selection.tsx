import React from "react";
import { cn } from "../../utils/cn";
import { useSharedState } from "./shared-state-provider";

interface ModelSelectionProps {
  className?: string;
}

const availableModels = [
  {
    id: "gemini-2.5-pro",
    name: "Gemini 2.5 Pro",
    provider: "Google",
    type: "Large",
  },
  {
    id: "claude-3-5-sonnet",
    name: "Claude 3.5 Sonnet",
    provider: "Anthropic",
    type: "Large",
  },
  { id: "llama-3-70b", name: "Llama 3 70B", provider: "Meta", type: "Large" },
  { id: "gpt-4o", name: "GPT-4o", provider: "OpenAI", type: "Large" },
  {
    id: "gemini-2.5-flash",
    name: "Gemini 2.5 Flash",
    provider: "Google",
    type: "Fast",
  },
  {
    id: "claude-3-haiku",
    name: "Claude 3 Haiku",
    provider: "Anthropic",
    type: "Fast",
  },
];

export const ModelSelection: React.FC<ModelSelectionProps> = ({
  className,
}) => {
  const { state, updateState } = useSharedState();

  const handleModelChange = (modelId: string) => {
    updateState({
      activeModelId: modelId,
    });
  };

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 shadow-md p-4",
        className,
      )}
    >
      <h3 className="text-lg font-medium mb-2">Model Selection</h3>
      <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
        Select the AI model to power your agent
      </p>

      <div className="space-y-2">
        {availableModels.map((model) => (
          <div
            key={model.id}
            className={cn(
              "p-3 border rounded-md cursor-pointer transition-colors",
              state.activeModelId === model.id
                ? "border-emerald-500 bg-emerald-500/10"
                : "border-gray-300 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-600",
            )}
            onClick={() => handleModelChange(model.id)}
          >
            <div className="flex justify-between items-center">
              <div>
                <div className="font-medium text-gray-900 dark:text-white">
                  {model.name}
                </div>
                <div className="text-sm text-gray-500 dark:text-gray-400">
                  {model.provider} • {model.type}
                </div>
              </div>
              {state.activeModelId === model.id && (
                <div className="h-4 w-4 rounded-full bg-emerald-500"></div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ModelSelection;
