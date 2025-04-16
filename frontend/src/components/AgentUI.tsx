import React, { useEffect, useState } from 'react';
import { useAgentState } from '../hooks/useAgentState';
import { useComponentStream } from '../hooks/useComponentStream';
import { Button } from './ui/Button';
import { Card } from './ui/Card';
import { CodeBlock } from './ui/CodeBlock';
import { ProgressBar } from './ui/ProgressBar';
import { Terminal } from './ui/Terminal';
import { Markdown } from './ui/Markdown';
import { AgentAction } from './AgentAction';
import { ToolOutput } from './ToolOutput';

interface AgentUIProps {
  agentId: string;
}

export const AgentUI: React.FC<AgentUIProps> = ({ agentId }) => {
  const { agentState } = useAgentState(agentId);
  const { components } = useComponentStream(agentId);
  const [activeView, setActiveView] = useState<'actions' | 'tools' | 'timeline'>('actions');

  const filteredComponents = components.filter(component => {
    if (activeView === 'actions') {
      return component.action_id !== undefined;
    } else if (activeView === 'tools') {
      return component.tool_id !== undefined;
    }
    return true; // Timeline shows all
  });

  return (
    <div className="flex flex-col h-full bg-charcoal-900 text-gray-100">
      <div className="flex justify-between items-center p-4 border-b border-gray-700">
        <h2 className="text-xl font-semibold">
          {agentState?.name || `Agent ${agentId}`}
        </h2>
        <div className="flex space-x-2">
          <Button 
            variant={activeView === 'actions' ? 'primary' : 'secondary'}
            onClick={() => setActiveView('actions')}
          >
            Actions
          </Button>
          <Button 
            variant={activeView === 'tools' ? 'primary' : 'secondary'}
            onClick={() => setActiveView('tools')}
          >
            Tools
          </Button>
          <Button 
            variant={activeView === 'timeline' ? 'primary' : 'secondary'}
            onClick={() => setActiveView('timeline')}
          >
            Timeline
          </Button>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {filteredComponents.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <p className="text-gray-400">No components to display</p>
          </div>
        ) : (
          filteredComponents.map(component => (
            <ComponentRenderer key={component.id} component={component} />
          ))
        )}
      </div>
      
      <div className="p-4 border-t border-gray-700">
        <ProgressBar 
          progress={agentState?.progress || 0} 
          status={agentState?.status || 'idle'} 
        />
      </div>
    </div>
  );
};

interface ComponentRendererProps {
  component: any;
}

const ComponentRenderer: React.FC<ComponentRendererProps> = ({ component }) => {
  switch (component.type) {
    case 'button':
      return <Button {...component.props} />;
    
    case 'card':
      return <Card {...component.props} />;
    
    case 'codeblock':
      return <CodeBlock {...component.props} />;
    
    case 'terminal':
      return <Terminal {...component.props} />;
    
    case 'markdown':
      return <Markdown {...component.props} />;
    
    case 'progress':
      return <ProgressBar {...component.props} />;
    
    case 'agentaction':
      return <AgentAction {...component.props} />;
    
    case 'tooloutput':
      return <ToolOutput {...component.props} />;
    
    default:
      return (
        <Card>
          <h3 className="text-lg font-medium">{component.type}</h3>
          <pre className="mt-2 p-2 bg-gray-800 rounded text-sm overflow-x-auto">
            {JSON.stringify(component.props, null, 2)}
          </pre>
        </Card>
      );
  }
};

export default AgentUI;
