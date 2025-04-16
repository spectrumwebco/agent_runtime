import React from 'react';
import { useSharedState } from '../../../context/shared-state-context';
import { cn } from '../../../utils/cn';

interface ProgressTrackerProps {
  className?: string;
}

export const ProgressTracker: React.FC<ProgressTrackerProps> = ({ className }) => {
  const { state } = useSharedState();
  const { progress, status, completedTasks, pendingTasks } = state.agentState;

  return (
    <div className={cn('flex flex-col space-y-4', className)}>
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-medium">Agent Progress</h3>
        <span className="text-sm px-2 py-1 rounded-full bg-emerald-500/20 text-emerald-500">
          {status === 'idle' && 'Idle'}
          {status === 'running' && 'Running'}
          {status === 'completed' && 'Completed'}
          {status === 'error' && 'Error'}
        </span>
      </div>
      
      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
        <div 
          className="bg-emerald-500 h-2.5 rounded-full transition-all duration-500 ease-in-out"
          style={{ width: `${progress}%` }}
        />
      </div>
      
      <div className="text-sm text-gray-500 dark:text-gray-400">
        {progress}% complete
      </div>
      
      {completedTasks.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium mb-2">Completed Tasks</h4>
          <ul className="space-y-1">
            {completedTasks.map((task, index) => (
              <li key={index} className="flex items-center text-sm">
                <span className="mr-2 text-emerald-500">✓</span>
                {task}
              </li>
            ))}
          </ul>
        </div>
      )}
      
      {pendingTasks.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium mb-2">Pending Tasks</h4>
          <ul className="space-y-1">
            {pendingTasks.map((task, index) => (
              <li key={index} className="flex items-center text-sm">
                <span className="mr-2 text-gray-400">○</span>
                {task}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default ProgressTracker;
