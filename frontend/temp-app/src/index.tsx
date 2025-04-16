import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';

const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-center">
      <h1 className="text-4xl font-bold text-emerald-500 mb-4">Kled.io Frontend</h1>
      <p className="text-xl">Agent Runtime Implementation</p>
    </div>
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
