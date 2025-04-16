import React from 'react';
import { Card } from './ui/Card';
import { CodeBlock } from './ui/CodeBlock';
import { Markdown } from './ui/Markdown';

interface AgentActionProps {
  agentID: string;
  actionID: string;
  actionType: string;
  thought?: string;
  code?: string;
  language?: string;
  result?: string;
  error?: string;
  timestamp?: number;
}

export const AgentAction: React.FC<AgentActionProps> = ({
  agentID,
  actionID,
  actionType,
  thought,
  code,
  language = 'javascript',
  result,
  error,
  timestamp,
}) => {
  const formattedTime = timestamp 
    ? new Date(timestamp * 1000).toLocaleTimeString() 
    : '';

  const renderActionContent = () => {
    switch (actionType) {
      case 'thinking':
        return thought ? <Markdown content={thought} /> : null;
      
      case 'code_generation':
        return (
          <>
            {thought && <Markdown content={thought} className="mb-4" />}
            {code && <CodeBlock code={code} language={language} />}
          </>
        );
      
      case 'error':
        return (
          <div className="p-4 bg-red-900/30 border border-red-700 rounded-md text-red-200">
            <h4 className="font-medium mb-2">Error</h4>
            <p>{error}</p>
          </div>
        );
      
      case 'success':
        return (
          <div className="p-4 bg-emerald-900/30 border border-emerald-700 rounded-md text-emerald-200">
            <h4 className="font-medium mb-2">Success</h4>
            <p>{result}</p>
          </div>
        );
      
      default:
        return (
          <>
            {thought && <Markdown content={thought} />}
            {result && (
              <div className="mt-4 p-3 bg-gray-800 rounded-md">
                <pre className="whitespace-pre-wrap">{result}</pre>
              </div>
            )}
          </>
        );
    }
  };

  return (
    <Card className="border-l-4 border-l-emerald-500">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-medium capitalize">
          {actionType.replace('_', ' ')}
        </h3>
        {formattedTime && (
          <span className="text-sm text-gray-400">{formattedTime}</span>
        )}
      </div>
      
      <div className="mt-2">
        {renderActionContent()}
      </div>
    </Card>
  );
};
