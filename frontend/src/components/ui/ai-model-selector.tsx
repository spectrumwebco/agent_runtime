import React from "react";
import { cn } from "../../utils/cn";

export interface AIModel {
  id: string;
  name: string;
  provider: string;
  description: string;
  isVerified: boolean;
  capabilities: string[];
}

interface AIModelSelectorProps {
  models: AIModel[];
  selectedModelId: string;
  onSelectModel: (modelId: string) => void;
  className?: string;
}

export const AIModelSelector: React.FC<AIModelSelectorProps> = ({
  models,
  selectedModelId,
  onSelectModel,
  className,
}) => {
  return (
    <div className={cn("flex flex-col space-y-2", className)}>
      <h3 className="text-lg font-medium mb-2">AI Model Selection</h3>
      <div className="grid grid-cols-1 gap-3">
        {models.map((model) => (
          <button
            key={model.id}
            onClick={() => onSelectModel(model.id)}
            className={cn(
              "flex items-start p-4 rounded-lg border transition-all",
              selectedModelId === model.id
                ? "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20"
                : "border-gray-200 dark:border-gray-700 hover:border-emerald-500"
            )}
          >
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <h4 className="font-medium">{model.name}</h4>
                {model.isVerified && (
                  <span className="inline-flex items-center rounded-full bg-emerald-100 dark:bg-emerald-900/30 px-2 py-0.5 text-xs font-medium text-emerald-800 dark:text-emerald-300">
                    Verified
                  </span>
                )}
              </div>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                {model.description}
              </p>
              <div className="flex flex-wrap gap-1 mt-2">
                {model.capabilities.map((capability, index) => (
                  <span
                    key={index}
                    className="inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-800 px-2 py-0.5 text-xs font-medium text-gray-800 dark:text-gray-300"
                  >
                    {capability}
                  </span>
                ))}
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};

export default AIModelSelector;
