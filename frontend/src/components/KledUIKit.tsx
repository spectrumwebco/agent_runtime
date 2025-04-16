import React from 'react';
import { SharedStateProvider } from '../context/shared-state-context';
import { FlippedLayout } from './layout/flipped-layout/flipped-layout';
import { ProgressTracker } from './agent/tracking/progress-tracker';
import { StateVisualizer } from './agent/visualization/state-visualizer';
import { ModelSelector } from './ui/model-selector/model-selector';
import { AIChat } from './ui/vercel-ai/ai-chat';
import { ThreeDCard } from './ui/aceternity/3d-card';
import { GradientCard } from './ui/aceternity/gradient-card';

interface KledUIKitProps {
  children?: React.ReactNode;
}

/**
 * KledUIKit - Main UI Kit component that integrates all libraries
 * 
 * This component provides a comprehensive UI framework for building
 * AI Agent interfaces with a flipped layout (IDE on left, controls on right)
 */
export const KledUIKit: React.FC<KledUIKitProps> = ({ children }) => {
  const leftPanelContent = (
    <div className="h-full w-full bg-gray-900 text-white p-4">
      <h2 className="text-xl font-bold mb-4">Code Editor</h2>
      <div className="bg-gray-800 p-4 rounded-lg h-[calc(100%-2rem)] overflow-auto">
        <pre className="text-sm">
          <code>{`// Example code
function runAgent() {
  const agent = new KledAgent({
    model: "gemini-2.5-pro",
    tools: defaultTools,
  });
  
  return agent.execute({
    task: "Build a React component",
    context: "Frontend development",
  });
}

`}</code>
        </pre>
      </div>
    </div>
  );

  const rightPanelContent = (
    <div className="h-full w-full p-4 space-y-6 overflow-auto">
      <h2 className="text-xl font-bold">Agent Controls</h2>
      
      <ThreeDCard className="p-4 bg-gray-800 dark:bg-gray-800">
        <ProgressTracker />
      </ThreeDCard>
      
      <GradientCard className="p-4">
        <StateVisualizer />
      </GradientCard>
      
      <ModelSelector />
      
      <div className="h-80">
        <AIChat />
      </div>
    </div>
  );

  return (
    <SharedStateProvider>
      <div className="h-screen w-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <header className="h-14 border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 flex items-center px-4">
          <h1 className="text-xl font-bold text-emerald-500">Kled.io Agent UI</h1>
        </header>
        
        <main className="h-[calc(100vh-3.5rem)]">
          <FlippedLayout
            leftPanel={leftPanelContent}
            rightPanel={rightPanelContent}
          />
        </main>
        
        {children}
      </div>
    </SharedStateProvider>
  );
};

export default KledUIKit;
