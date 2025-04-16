import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { RootProvider } from './providers/root-provider';
import { SharedStateProvider } from './components/ui/shared-state-provider';
import KledMainLayout from './components/ui/kled-main-layout';
import { CoPilotIntegration } from './components/ui/copilot-integration';

const App = () => {
  return (
    <RootProvider>
      <SharedStateProvider>
        <KledMainLayout>
          <div className="grid grid-cols-1 gap-4 p-4">
            <CoPilotIntegration />
          </div>
        </KledMainLayout>
      </SharedStateProvider>
    </RootProvider>
  );
};

const root = document.getElementById('root');
if (root) {
  ReactDOM.createRoot(root).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
}
