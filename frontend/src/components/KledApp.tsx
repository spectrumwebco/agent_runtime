import React from "react";
import { ThemeProvider } from "./ui/theme-provider";
import { SharedStateProvider } from "./ui/shared-state-provider";
import { KledMainLayout } from "./ui/kled-main-layout";
import { AIChat } from "./ui/vercel-ai/ai-chat";

export const KledApp: React.FC = () => {
  return (
    <ThemeProvider defaultTheme="dark">
      <SharedStateProvider>
        <KledMainLayout>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 h-full">
            <div className="h-full overflow-hidden border rounded-lg">
              <div className="h-full bg-white dark:bg-gray-800 p-4">
                <h2 className="text-xl font-bold mb-4">Code Editor</h2>
                <div className="h-[calc(100%-2rem)] bg-gray-100 dark:bg-gray-900 rounded-lg p-4">
                  <pre className="text-sm">
                    <code>
                      {`// Kled.io Code Editor
import { useState } from 'react';

function KledComponent() {
  const [state, setState] = useState({});
  
  return (
    <div>
      <h1>Kled.io Component</h1>
      <p>Edit this code to get started</p>
    </div>
  );
}`}
                    </code>
                  </pre>
                </div>
              </div>
            </div>
            
            <div className="h-full overflow-hidden border rounded-lg">
              <AIChat />
            </div>
          </div>
        </KledMainLayout>
      </SharedStateProvider>
    </ThemeProvider>
  );
};

export default KledApp;
