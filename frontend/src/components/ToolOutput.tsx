import React from 'react';
import { Card } from './ui/Card';
import { CodeBlock } from './ui/CodeBlock';
import { Markdown } from './ui/Markdown';

interface ToolOutputProps {
  agentID: string;
  toolID: string;
  toolName: string;
  toolInput: Record<string, any>;
  toolOutput: Record<string, any>;
  timestamp?: number;
}

export const ToolOutput: React.FC<ToolOutputProps> = ({
  agentID,
  toolID,
  toolName,
  toolInput,
  toolOutput,
  timestamp,
}) => {
  const formattedTime = timestamp 
    ? new Date(timestamp * 1000).toLocaleTimeString() 
    : '';

  const displayToolName = toolName
    .replace(/([A-Z])/g, ' $1') // Add spaces before capital letters
    .replace(/_/g, ' ') // Replace underscores with spaces
    .replace(/^./, str => str.toUpperCase()); // Capitalize first letter

  const renderToolInput = () => {
    if (typeof toolInput === 'string') {
      return <p className="text-gray-300">{toolInput}</p>;
    }
    
    if (toolInput.code) {
      return (
        <CodeBlock 
          code={toolInput.code} 
          language={toolInput.language || 'javascript'} 
        />
      );
    }
    
    if (toolInput.query) {
      return <p className="text-gray-300">{toolInput.query}</p>;
    }
    
    return (
      <div className="bg-gray-800 p-3 rounded-md">
        <pre className="text-sm overflow-x-auto whitespace-pre-wrap">
          {JSON.stringify(toolInput, null, 2)}
        </pre>
      </div>
    );
  };

  const renderToolOutput = () => {
    if (typeof toolOutput === 'string') {
      return <p className="text-gray-200">{toolOutput}</p>;
    }
    
    if (toolOutput.markdown) {
      return <Markdown content={toolOutput.markdown} />;
    }
    
    if (toolOutput.code) {
      return (
        <CodeBlock 
          code={toolOutput.code} 
          language={toolOutput.language || 'javascript'} 
        />
      );
    }
    
    if (toolOutput.result) {
      return <p className="text-gray-200">{toolOutput.result}</p>;
    }
    
    return (
      <div className="bg-gray-800 p-3 rounded-md">
        <pre className="text-sm overflow-x-auto whitespace-pre-wrap">
          {JSON.stringify(toolOutput, null, 2)}
        </pre>
      </div>
    );
  };

  return (
    <Card className="border-l-4 border-l-blue-500">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-medium">
          {displayToolName}
        </h3>
        {formattedTime && (
          <span className="text-sm text-gray-400">{formattedTime}</span>
        )}
      </div>
      
      <div className="mt-4 space-y-4">
        <div>
          <h4 className="text-sm font-medium text-gray-400 mb-2">Input</h4>
          {renderToolInput()}
        </div>
        
        <div>
          <h4 className="text-sm font-medium text-gray-400 mb-2">Output</h4>
          {renderToolOutput()}
        </div>
      </div>
    </Card>
  );
};
