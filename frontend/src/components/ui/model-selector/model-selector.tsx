import React from 'react';
import { useSharedState } from '../../../context/shared-state-context';
import { cn } from '../../../utils/cn';

interface ModelSelectorProps {
  className?: string;
}

export const ModelSelector: React.FC<ModelSelectorProps> = ({ className }) => {
  const { state, updateModelState } = useSharedState();
  const { selectedModel, availableModels } = state.modelState;

  const modelInfo = {
    'gemini-2.5-pro': {
      name: 'Gemini 2.5 Pro',
      description: 'Specialized for coding tasks',
      icon: '💻',
    },
    'llama-4-scout': {
      name: 'Llama 4 Scout',
      description: 'Optimized for standard operations',
      icon: '🔍',
    },
    'llama-4-maverick': {
      name: 'Llama 4 Maverick',
      description: 'Specialized for reasoning tasks',
      icon: '🧠',
    },
    'gpt-4o': {
      name: 'GPT-4o',
      description: 'General purpose model',
      icon: '🤖',
    },
  };

  const handleModelChange = (model: string) => {
    updateModelState({ selectedModel: model });
  };

  return (
    <div className={cn('flex flex-col space-y-4', className)}>
      <h3 className="text-lg font-medium">AI Model Selection</h3>
      
      <div className="grid grid-cols-1 gap-3">
        {availableModels.map((model) => {
          const info = modelInfo[model as keyof typeof modelInfo] || {
            name: model,
            description: 'AI model',
            icon: '🤖',
          };
          
          return (
            <button
              key={model}
              className={cn(
                'flex items-center p-3 rounded-lg border-2 transition-all duration-200',
                selectedModel === model
                  ? 'border-emerald-500 bg-emerald-500/10'
                  : 'border-gray-200 dark:border-gray-700 hover:border-emerald-500/50'
              )}
              onClick={() => handleModelChange(model)}
            >
              <div className="text-2xl mr-3">{info.icon}</div>
              <div className="flex flex-col items-start">
                <span className="font-medium">{info.name}</span>
                <span className="text-xs text-gray-500 dark:text-gray-400">
                  {info.description}
                </span>
              </div>
              {selectedModel === model && (
                <div className="ml-auto text-emerald-500">✓</div>
              )}
            </button>
          );
        })}
      </div>
      
      <div className="text-xs text-gray-500 dark:text-gray-400 mt-2">
        Different models are optimized for specific tasks in the agent system.
      </div>
    </div>
  );
};

export default ModelSelector;
