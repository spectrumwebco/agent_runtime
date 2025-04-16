import React, { useState } from "react";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";
import { useSharedState } from "./shared-state-provider";

interface Model {
  id: string;
  name: string;
  type: "reasoning" | "coding" | "standard";
  description: string;
  provider: string;
}

interface ModelSelectionProps {
  className?: string;
}

export const ModelSelection: React.FC<ModelSelectionProps> = ({
  className,
}) => {
  const { state, updateState } = useSharedState();
  const [isChanging, setIsChanging] = useState(false);
  
  const models: Model[] = [
    {
      id: "gemini-2.5-pro",
      name: "Gemini 2.5 Pro",
      type: "coding",
      description: "Specialized for coding tasks with strong code generation capabilities",
      provider: "Google",
    },
    {
      id: "llama-4-scout",
      name: "Llama 4 Scout",
      type: "standard",
      description: "Optimized for standard operations with balanced capabilities",
      provider: "Meta",
    },
    {
      id: "llama-4-maverick",
      name: "Llama 4 Maverick",
      type: "reasoning",
      description: "Enhanced reasoning capabilities for complex problem-solving",
      provider: "Meta",
    },
  ];

  const handleModelChange = (modelId: string) => {
    setIsChanging(true);
    
    setTimeout(() => {
      updateState({ activeModelId: modelId });
      setIsChanging(false);
    }, 1000);
  };

  const activeModelId = state.activeModelId || "llama-4-scout";

  return (
    <SpotlightCard className={cn("p-4", className)}>
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-medium">Model Selection</h3>
        <span className={cn(
          "px-2 py-1 rounded-full text-xs",
          isChanging ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300" : "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300"
        )}>
          {isChanging ? "Changing..." : "Ready"}
        </span>
      </div>

      <div className="space-y-4">
        {models.map((model) => (
          <div
            key={model.id}
            className={cn(
              "border rounded-lg p-4 transition-colors cursor-pointer hover:border-emerald-500",
              model.id === activeModelId ? "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/10" : "border-gray-200 dark:border-gray-700"
            )}
            onClick={() => !isChanging && handleModelChange(model.id)}
          >
            <div className="flex justify-between items-start mb-2">
              <div>
                <h4 className="font-medium">{model.name}</h4>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {model.provider}
                </p>
              </div>
              <span className={cn(
                "px-2 py-1 rounded-full text-xs",
                model.type === "coding" ? "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300" :
                model.type === "reasoning" ? "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300" :
                "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300"
              )}>
                {model.type}
              </span>
            </div>
            
            <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
              {model.description}
            </p>
            
            {model.id === activeModelId ? (
              <GradientButton
                size="sm"
                className="w-full"
                disabled
              >
                Currently Active
              </GradientButton>
            ) : (
              <GradientButton
                variant="outline"
                size="sm"
                className="w-full"
                onClick={() => !isChanging && handleModelChange(model.id)}
                disabled={isChanging}
              >
                Switch to this Model
              </GradientButton>
            )}
          </div>
        ))}
      </div>
    </SpotlightCard>
  );
};

export default ModelSelection;
