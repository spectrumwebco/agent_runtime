import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { SharedStateProvider } from './components/ui/shared-state-provider';
import KledMainLayout from './components/ui/kled-main-layout';
import { CoPilotIntegration } from './components/ui/copilot-integration';

const App = () => {
  return (
    <SharedStateProvider>
      <KledMainLayout>
        <div className="grid grid-cols-1 gap-4">
          <CoPilotIntegration />
        </div>
      </KledMainLayout>
    </SharedStateProvider>
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
