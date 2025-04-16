import React from 'react';
import { useSharedState } from '../../../context/shared-state-context';
import { cn } from '../../../utils/cn';

interface StateVisualizerProps {
  className?: string;
}

export const StateVisualizer: React.FC<StateVisualizerProps> = ({ className }) => {
  const { state } = useSharedState();
  const { status, currentTask } = state.agentState;

  const states = [
    { id: 'idle', label: 'Idle', color: 'bg-gray-400' },
    { id: 'planning', label: 'Planning', color: 'bg-blue-500' },
    { id: 'executing', label: 'Executing', color: 'bg-yellow-500' },
    { id: 'verifying', label: 'Verifying', color: 'bg-purple-500' },
    { id: 'completed', label: 'Completed', color: 'bg-emerald-500' },
    { id: 'error', label: 'Error', color: 'bg-red-500' },
  ];

  const activeStateIndex = states.findIndex(s => 
    (status === 'running' && currentTask.toLowerCase().includes(s.id)) || 
    s.id === status
  );

  return (
    <div className={cn('flex flex-col space-y-4', className)}>
      <h3 className="text-lg font-medium">Agent State</h3>
      
      <div className="flex flex-col space-y-2">
        {states.map((state, index) => (
          <div 
            key={state.id}
            className={cn(
              'flex items-center p-3 rounded-lg transition-all duration-300',
              index === activeStateIndex 
                ? `${state.color}/20 border-l-4 ${state.color}` 
                : 'bg-gray-100 dark:bg-gray-800'
            )}
          >
            <div 
              className={cn(
                'w-3 h-3 rounded-full mr-3',
                index === activeStateIndex ? state.color : 'bg-gray-300 dark:bg-gray-600'
              )}
            />
            <span className={cn(
              'font-medium',
              index === activeStateIndex ? 'text-gray-900 dark:text-white' : 'text-gray-500 dark:text-gray-400'
            )}>
              {state.label}
            </span>
          </div>
        ))}
      </div>
      
      {currentTask && (
        <div className="mt-2 text-sm text-gray-500 dark:text-gray-400">
          Current task: {currentTask}
        </div>
      )}
    </div>
  );
};

export default StateVisualizer;
